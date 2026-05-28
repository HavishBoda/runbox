package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"
	// "sync"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/joho/godotenv"
)

func runCode(code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: "python:3.11-alpine",
			Cmd:   []string{"python3", "-c", code},
		},
		HostConfig: &container.HostConfig{
			Resources: container.Resources{
				Memory:   128 * 1024 * 1024,
				CPUQuota: 50000,
			},
			NetworkMode: "none",
		},
	})
	if err != nil {
		return "", err
	}

	defer cli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{Force: true})

	_, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	wait := cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})
	select {
	case err := <-wait.Error:
		if err != nil {
			return "", err
		}
	case <-wait.Result:
	}

	logsResult, err := cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", err
	}
	defer logsResult.Close()

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, logsResult)
	if err != nil {
		return "", err
	}

	return stdout.String() + stderr.String(), nil
}

func main() {
    godotenv.Load()

	code needed
    result, err := runAgent("Calculate the sum of all even numbers between 1 and 100")
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

	fmt.Println("Result:", result)

	// factual, no code needed
	result, err = runAgent("What is the capital of France")
	if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Result:", result)
}