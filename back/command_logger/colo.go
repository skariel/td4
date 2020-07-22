package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please provide a command to execute")
	}

	cmdname := strings.Replace(os.Args[1], "./", "", 1)

	f, err := os.Create("co_" + cmdname + "_" + uuid.New().String() + ".lo")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(os.Args[1])
	cmd.Stdout = f
	cmd.Stderr = f

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}
