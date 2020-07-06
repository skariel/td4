package main

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"time"

	"td4/back/db"
	gdb "td4/back/db/generated"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/joshdk/go-junit"
)

// configuration consts
const (
	sleepTime       = 5 * time.Second
	dockerImageName = "td4:v1"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// connect to the DB
	q, dbase, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	// main work loop: get run, run, report, repeat!
	shouldSleep := false

	for {
		if shouldSleep {
			time.Sleep(sleepTime)
		}

		runs, err := q.FetchSomeRun(context.Background())
		if err != nil {
			log.Printf("error while fetching run: %v", err)

			shouldSleep = true

			continue
		}

		if len(runs) == 0 {
			shouldSleep = true
			continue
		}

		shouldSleep = false

		// we have a run!
		run := runs[0]
		log.Printf("new run arrived id = %v", run.ID)

		sol, tes, conf, err := getCodesAndConf(ctx, &run, q)
		if err != nil {
			log.Printf("error getting run data: %v for run: %v", err, run)

			shouldSleep = true

			continue
		}

		suites, err := runContainer(ctx, cli, tes, sol, conf, q, dbase, run.ID)
		if err != nil {
			log.Printf("error running container: %v", err)

			shouldSleep = true

			continue
		}

		if suites == nil {
			log.Print("no results to show")

			continue
		}

		if len(suites) == 0 {
			log.Print("zero results")

			continue
		}
	}
}

func runContainer(
	ctx context.Context,
	cli *client.Client,
	tes *gdb.GetTestCodeByIDRow,
	sol *gdb.Td4SolutionCode,
	conf *gdb.Td4RunConfig,
	q *gdb.Queries,
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
	}, nil, "")

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
		id := resp.ID

		time.Sleep(time.Duration(conf.MaxTimeSecs) * time.Second)

		containers, errg := cli.ContainerList(ctx, types.ContainerListOptions{})
		if errg != nil {
			log.Printf("error listing containers: %v", err)
		}

		// only kill container if still listed...
		for ix := range containers {
			c := containers[ix]
			if c.ID == id {
				log.Printf("timeout! stopping container id=%v and ending run id=%v with status=stop", id, runid)

				err = q.EndRunByID(ctx, gdb.EndRunByIDParams{
					ID:     runid,
					Status: gdb.Td4TypeRunStatusStop})
				if err != nil {
					log.Printf("error reporting stop status for long running run: %v", err)
				}

				to := time.Duration(1) * time.Second

				err = cli.ContainerStop(ctx, id, &to)
				if err != nil {
					log.Printf("error stopping container after timeout: %v", err)
				}

				break
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
	q *gdb.Queries,
	dbase *sql.DB,
	runid int32) ([]junit.Suite, error) {
	// get test output!
	r, _, err := cli.CopyFromContainer(ctx, resp.ID, "/test_result.xml")

	if r != nil {
		defer func() { _ = r.Close() }()
	} else {
		return nil, errors.New("could not find reults xml file")
	}

	if err != nil {
		return nil, err
	}

	xml, err := readContents(r)

	if err != nil {
		return nil, err
	}

	if xml == nil {
		return nil, errors.New("could not read result xml contents")
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

	var status gdb.Td4TypeRunStatus

	if (suites[0].Totals.Failed == 0) && (suites[0].Totals.Error == 0) {
		status = gdb.Td4TypeRunStatusPass
	} else {
		status = gdb.Td4TypeRunStatusFail
	}

	// checking error at end of transaction
	_ = tq.EndRunByID(ctx, gdb.EndRunByIDParams{
		ID:     runid,
		Status: status})

	for ix := range suites {
		suite := suites[ix]
		for _, test := range suite.Tests {
			if test.Error != nil {
				_, _ = tq.InsertRunResult(ctx, gdb.InsertRunResultParams{
					RunID:  runid,
					Status: gdb.Td4TypeRunResultStatusFail,
					Title:  sql.NullString{String: test.Name, Valid: true},
					Output: sql.NullString{String: test.Error.Error(), Valid: true}})
			} else {
				_, _ = tq.InsertRunResult(ctx, gdb.InsertRunResultParams{
					RunID:  runid,
					Status: gdb.Td4TypeRunResultStatusPass,
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
	tes *gdb.GetTestCodeByIDRow,
	sol *gdb.Td4SolutionCode,
	_ *gdb.Td4RunConfig) error {
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

	_, err = cli.ContainerWait(ctx, resp.ID)

	return err
}

func getCodesAndConf(
	ctx context.Context,
	run *gdb.Td4Run,
	q *gdb.Queries,
) (*gdb.Td4SolutionCode, *gdb.GetTestCodeByIDRow, *gdb.Td4RunConfig, error) {
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
