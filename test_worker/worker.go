package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"log"
	"os"
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

	log.Println("Cretaing a Docker client")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	log.Println("pulling docker python image")

	reader, err := cli.ImagePull(ctx, "docker.io/library/python", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Println("DOCKER: " + scanner.Text())

		if err = scanner.Err(); err != nil {
			log.Printf("error reading output from docker image pull: %v", err)
		}
	}

	// connect to the DB
	q, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	for {
		log.Println(".")
		time.Sleep(sleepTimeSeconds * time.Second)

		runs, err := q.FetchSomeRun(context.Background())
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		if len(runs) == 0 {
			continue
		}

		// we have a run!
		run := runs[0]
		log.Printf("new run arrived: %v", run)

		// get the config, test and solution codes
		sol, err := q.GetSolutionByID(context.Background(), run.SolutionCodeID)
		if err != nil {
			log.Printf("error getting solution code: %v for run: %v", err, run)
			continue
		}

		tes, err := q.GetTestCodeByID(context.Background(), sol.TestCodeID)
		if err != nil {
			log.Printf("error getting test code: %v for solution: %v", err, sol)
			continue
		}

		conf, err := q.GetConfByDiplayName(context.Background(), run.RunConfig)
		if err != nil {
			log.Printf("error getting conf: %v for run: %v", err, run)
			continue
		}

		log.Println("Creating new container")

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "python",
			Cmd:   []string{"python", "test.py"},
		}, nil, nil, "")
		if err != nil {
			log.Printf("error creating python container: %v", err)
			continue
		}

		// TODO: copy files to a docker container
		var buf bytes.Buffer

		content := []byte(tes.Code)
		tw := tar.NewWriter(&buf)
		err = tw.WriteHeader(&tar.Header{
			Name: "test.py",           // filename
			Mode: 0777,                // permissions
			Size: int64(len(content)), // filesize
		})

		if err != nil {
			log.Printf("docker copy header test code: %v", err)
			continue
		}

		_n, err := tw.Write(content)
		if err != nil {
			log.Printf("docker write content test code: %v", err)
			continue
		}

		if _n != len(content) {
			log.Printf("docker could not write all bytes content test code: %v, != %v", len(content), _n)
			continue
		}

		err = tw.Close()
		if err != nil {
			log.Printf("docker closinf file test code: %v", err)
			continue
		}

		// use &buf as argument for content in CopyToContainer
		err = cli.CopyToContainer(ctx, resp.ID, "/", &buf, types.CopyToContainerOptions{})
		if err != nil {
			log.Printf("docker copy test code: %v", err)
			continue
		}

		// TODO: run the docker container

		log.Println("running the container")

		if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("error running the container: %v", err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err = <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)

		log.Printf("code: %v", tes.Code)
		log.Printf("solution: %v", sol.Code)
		log.Printf("Conf: %v", conf)
		// TODO: check if timeout
		// TODO: check logs / result
		// TODO: update run results
	}
}
