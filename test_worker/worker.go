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

	log.Println("cretaing a Docker client")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)

	log.Println("checking if image python exists")

	imageSummaries, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Fatalf("cannot check docker images: %v", err)
	}

	var pythonImageExists = false

	for i := range imageSummaries {
		if strings.HasPrefix(imageSummaries[i].RepoTags[0], "python") {
			log.Println("found docker python image")

			pythonImageExists = true

			break
		}
	}

	if !pythonImageExists {
		log.Println("pulling docker python image")

		reader, _err := cli.ImagePull(ctx, "docker.io/library/python", types.ImagePullOptions{})
		if _err != nil {
			panic(_err)
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Println("DOCKER: " + scanner.Text())

			if err = scanner.Err(); err != nil {
				log.Printf("error reading output from docker image pull: %v", err)
			}
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

		// copy files to the docker container
		err = copyToDocker(cli, ctx, &resp, tes.Code, "test.py")
		if err != nil {
			log.Printf("error while copying test code to docker: %v", err)
			continue
		}

		err = copyToDocker(cli, ctx, &resp, sol.Code, "solution.py")
		if err != nil {
			log.Printf("error while copying solution code to docker: %v", err)
			continue
		}

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

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			panic(err)
		}

		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		if err != nil {
			log.Printf("error while running container: %v", err)
			continue
		}

		log.Printf("code: %v", tes.Code)
		log.Printf("solution: %v", sol.Code)
		log.Printf("Conf: %v", conf)
		// TODO: check if timeout
		// TODO: check logs / result
		// TODO: update run results
	}
}

func copyToDocker(cli *client.Client, ctx context.Context, resp *container.ContainerCreateCreatedBody, data string, fn string) error {
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
