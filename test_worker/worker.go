package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"../sql/db"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// configuration consts
const (
	sleepTimeSeconds = 5
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	if !checkIfPythonImageExists(ctx, cli) {
		log.Println("pulling docker python image")
		pullPythonImage(ctx, cli)
	} else {
		log.Println("found docker python image")
	}

	// connect to the DB
	q, err := db.ConnectDB()
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
		log.Printf("new run arrived: %v", run)

		sol, tes, conf, err := getCodesAndConf(ctx, &run, q)
		if err != nil {
			log.Printf("error getting run data: %v for run: %v", err, run)
			continue
		}

		log.Println("running")

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "python",
			Cmd:   []string{"python", "-u", "test.py"},
		}, nil, nil, "")
		if err != nil {
			log.Printf("error creating python container: %v", err)
			continue
		}

		if err := runRun(ctx, cli, &resp, tes, sol, conf); err != nil {
			log.Printf("error while running container: %v", err)
		}
	}
}

func runRun(
	ctx context.Context,
	cli *client.Client,
	resp *container.ContainerCreateCreatedBody,
	tes *db.Td4TestCode,
	sol *db.Td4SolutionCode,
	_ *db.Td4RunConfig) error {
	// copy files to the docker container
	err := copyToDocker(ctx, cli, resp, tes.Code, "test.py")
	if err != nil {
		return err
	}

	err = copyToDocker(ctx, cli, resp, sol.Code, "solution.py")
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

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("error while reading container logs: %v", err)
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		return err
	}

	return nil
	// TODO: check if timeout
	// TODO: check logs / result
	// TODO: update run results
}

func getCodesAndConf(ctx context.Context, run *db.Td4Run, q *db.Queries) (*db.Td4SolutionCode, *db.Td4TestCode, *db.Td4RunConfig, error) {
	sol, err := q.GetSolutionByID(ctx, run.SolutionCodeID)
	if err != nil {
		return nil, nil, nil, err
	}

	tes, err := q.GetTestCodeByID(ctx, sol.TestCodeID)
	if err != nil {
		return nil, nil, nil, err
	}

	conf, err := q.GetConfByDiplayName(ctx, run.RunConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	return &sol, &tes, &conf, nil
}

func checkIfPythonImageExists(ctx context.Context, cli *client.Client) bool {
	log.Println("checking if image python exists")

	imageSummaries, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Fatalf("cannot check docker images: %v", err)
	}

	for i := range imageSummaries {
		if strings.HasPrefix(imageSummaries[i].RepoTags[0], "python") {
			return true
		}
	}

	return false
}

func pullPythonImage(ctx context.Context, cli *client.Client) {
	reader, err := cli.ImagePull(ctx, "docker.io/library/python", types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("error while pulling python image: %v", err)
		panic(err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Println("DOCKER: " + scanner.Text())

		if err = scanner.Err(); err != nil {
			log.Printf("error reading output from docker image pull: %v", err)
		}
	}
}

func copyToDocker(ctx context.Context, cli *client.Client, resp *container.ContainerCreateCreatedBody, data, fn string) error {
	var buf bytes.Buffer

	const perm = 0777

	content := []byte(data)
	tw := tar.NewWriter(&buf)
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

	err = tw.Close()
	if err != nil {
		return err
	}

	err = cli.CopyToContainer(ctx, resp.ID, "/", &buf, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	return nil
}
