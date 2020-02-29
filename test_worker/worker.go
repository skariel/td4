package main

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"../sql/db"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joshdk/go-junit"
)

// configuration consts
const (
	sleepTimeSeconds = 5
	dockerImageName  = "td4:v1"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	// connect to the DB
	q, dbase, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	// main work loop: get run, get run data, runRun, repeat!
	for {
		time.Sleep(sleepTimeSeconds * time.Second)

		runs, err := q.FetchSomeRun(context.Background())
		if err != nil {
			log.Printf("error while fetching run: %v", err)
			continue
		}

		if len(runs) == 0 {
			continue
		}

		// we have a run!
		run := runs[0]
		log.Printf("new run arrived id = %v", run.ID)

		sol, tes, conf, err := getCodesAndConf(ctx, &run, q)
		if err != nil {
			log.Printf("error getting run data: %v for run: %v", err, run)
			continue
		}

		suites, err := runContainer(ctx, cli, tes, sol, conf, q, dbase, run.ID)
		if err != nil {
			log.Printf("error running container: %v", err)
			continue
		}

		log.Print("Test results:")

		for ix := range suites {
			suite := suites[ix]
			fmt.Println(suite.Name)

			for _, test := range suite.Tests {
				fmt.Printf("  %s\n", test.Name)

				if test.Error != nil {
					fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
				} else {
					fmt.Printf("    %s\n", test.Status)
				}
			}
		}
	}
}

func runContainer(
	ctx context.Context,
	cli *client.Client,
	tes *db.Td4TestCode,
	sol *db.Td4SolutionCode,
	conf *db.Td4RunConfig,
	q *db.Queries,
	dbase *sql.DB,
	runid int32) ([]junit.Suite, error) {
	const mega = 1024 * 1024

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: dockerImageName,
		Cmd:   []string{"pytest", "--junitxml=./test_result.xml", "-rp", "./test/test.py"},
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:     int64(conf.MemoryUsageMb * mega),
			MemorySwap: int64(conf.MemoryUsageMb * mega),
			CPUPeriod:  int64(conf.CpuPeriod),
			CPUQuota:   int64(conf.CpuQuota),
		},
	},
		nil, "")

	if err != nil {
		return nil, err
	}

	defer func() {
		log.Print("removing container")

		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true})

		if err != nil {
			log.Printf("error while removing container: %v", err)
		}
	}()

	go func() {
		const one = 1 // ddont ask ;)

		id := resp.ID

		time.Sleep(time.Duration(conf.MaxTimeSecs) * time.Second)

		containers, errg := cli.ContainerList(ctx, types.ContainerListOptions{})
		if errg != nil {
			log.Printf("error listing containers: %v", err)
		}

		for ix := range containers {
			c := containers[ix]
			if c.ID == id {
				log.Printf("timeout! stopping container id=%v", id)

				to := time.Duration(one) * time.Second

				err = cli.ContainerStop(ctx, id, &to)
				if err != nil {
					log.Printf("error stopping container after timeout: %v", err)
				}
			}
		}
	}()

	err = runRun(ctx, cli, &resp, tes, sol, conf)
	if err != nil {
		return nil, err
	}

	suites, err := reportResults(ctx, cli, &resp, q, dbase, runid)
	if err != nil {
		return nil, err
	}

	return suites, nil
}

func reportResults(
	ctx context.Context,
	cli *client.Client,
	resp *container.ContainerCreateCreatedBody,
	q *db.Queries,
	dbase *sql.DB,
	runid int32) ([]junit.Suite, error) {
	// get test output!
	r, _, err := cli.CopyFromContainer(ctx, resp.ID, "/test_result.xml")

	if r != nil {
		defer func() { _ = r.Close() }()
	}

	if err != nil {
		return nil, err
	}

	xml, err := readContents(r)

	if err != nil {
		return nil, err
	}

	suites, err := junit.Ingest(xml)
	if err != nil {
		return nil, err
	}

	tx, err := dbase.Begin()
	if err != nil {
		return nil, err
	}

	tq := q.WithTx(tx)

	var status db.Td4TypeRunStatus

	if (suites[0].Totals.Failed == 0) && (suites[0].Totals.Error == 0) {
		status = db.Td4TypeRunStatusPass
	} else {
		status = db.Td4TypeRunStatusFail
	}

	// checking error at end of transaction
	_ = tq.UpdateRunStatusByID(ctx, db.UpdateRunStatusByIDParams{
		ID:     runid,
		Status: status})

	for ix := range suites {
		suite := suites[ix]
		for _, test := range suite.Tests {
			if test.Error != nil {
				_, _ = tq.InsertRunResult(ctx, db.InsertRunResultParams{
					RunID:  runid,
					Status: db.Td4TypeRunResultStatusFail,
					Title:  sql.NullString{String: test.Name, Valid: true},
					Output: sql.NullString{String: test.Error.Error(), Valid: true}})
			} else {
				_, _ = tq.InsertRunResult(ctx, db.InsertRunResultParams{
					RunID:  runid,
					Status: db.Td4TypeRunResultStatusPass,
					Title:  sql.NullString{String: test.Name, Valid: true}})
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return suites, nil
}

func readContents(reader io.Reader) ([]byte, error) {
	tr := tar.NewReader(reader)

	_, err := tr.Next()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(tr)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func runRun(
	ctx context.Context,
	cli *client.Client,
	resp *container.ContainerCreateCreatedBody,
	tes *db.Td4TestCode,
	sol *db.Td4SolutionCode,
	_ *db.Td4RunConfig) error {
	err := copyToDocker(ctx, cli, resp, tes.Code, "test/test.py")
	if err != nil {
		return err
	}

	err = copyToDocker(ctx, cli, resp, sol.Code, "test/solution.py")
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}

func getCodesAndConf(ctx context.Context, run *db.Td4Run, q *db.Queries) (*db.Td4SolutionCode, *db.Td4TestCode, *db.Td4RunConfig, error) {
	sol, err := q.RAWGetSolutionCodeByID(ctx, run.SolutionCodeID)
	if err != nil {
		return nil, nil, nil, err
	}

	tes, err := q.GetTestCodeByID(ctx, sol.TestCodeID)
	if err != nil {
		return nil, nil, nil, err
	}

	conf, err := q.GetConfByDisplayName(ctx, run.RunConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	return &sol, &tes, &conf, nil
}

func copyToDocker(ctx context.Context, cli *client.Client, resp *container.ContainerCreateCreatedBody, data, fn string) error {
	var buf bytes.Buffer

	const perm = 0777

	content := []byte(data)
	tw := tar.NewWriter(&buf)

	defer func() { _ = tw.Close() }()

	err := tw.WriteHeader(&tar.Header{
		Name: fn,                  // filename
		Mode: perm,                // permissions
		Size: int64(len(content)), // filesize
	})

	if err != nil {
		return err
	}

	_n, err := tw.Write(content)
	if err != nil {
		return err
	}

	if _n != len(content) {
		return errors.New("could not write all data in copyToDocker")
	}

	err = cli.CopyToContainer(ctx, resp.ID, "/", &buf, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}
