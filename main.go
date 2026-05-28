package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"
	"sync"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
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
    codes := []string{
        `print("execution 1")`,
        `print("execution 2")`,
        `print("execution 3")`,
    }

    var wg sync.WaitGroup

    for _, code := range codes {
        wg.Add(1)
        go func(c string) {
            defer wg.Done()
            output, err := runCode(c)
            if err != nil {
                fmt.Fprintf(os.Stderr, "error: %v\n", err)
                return
            }
            fmt.Print(output)
        }(code)
    }

    wg.Wait()
}