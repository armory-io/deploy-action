package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/armory-io/deploy-with-spinnaker-action/pkg/action"

	"github.com/armory-io/deploy-with-spinnaker-action/pkg/ops"

	"github.com/armory-io/deploy-with-spinnaker-action/pkg/gate"
)

func main() {
	ac, err := action.GetConfig()
	if err != nil {
		fmt.Printf("unable to execute action: %s\n", err.Error())
		os.Exit(1)
	}

	if err := ac.Validate(); err != nil {
		fmt.Printf("unable to execute action: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("beginning deployment of application %s to accounts %s\n", ac.Application, ac.Accounts)

	manifestProvider, err := providerFactory(ac.Manifest, ac.ManifestPath)
	if err != nil {
		fmt.Printf("manifests to be deployed are invalid: %s\n", err.Error())
		os.Exit(1)
	}

	tasks := []interface{}{}
	for _, account := range ac.Accounts {
		op, err := ops.NewKubernetesDeployTask(account, ac.Application, manifestProvider)
		if err != nil {
			fmt.Printf("failed to generate deployment task: %s\n", err.Error())
			os.Exit(1)
		}
		tasks = append(tasks, op)
	}

	task := ops.Task{
		Job:         tasks,
		Application: ac.Application,
		Description: "Deployment from Github Actions",
	}

	gateClient, err := gate.New(ac.Config, ac.BaseURL)
	if err != nil {
		fmt.Printf("failed to initialize gate client: %s\n", err.Error())
		os.Exit(1)
	}

	executionID, err := performDeploy(task, gateClient)
	if err != nil {
		fmt.Printf("failed to execute deployment task: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("deployment created with id %s\n", executionID)
	if ac.Wait != true {
		fmt.Println("deployment created successfully.")
		return
	}

	fmt.Println("waiting for deployment to complete...")
	timeout, _ := context.WithTimeout(context.Background(), time.Second*30)
	if err := pollExecution(timeout, executionID, gateClient); err != nil {
		fmt.Printf("there was an error waiting for deployment to complete: %s\n", err.Error())
		os.Exit(1)
	}

	execution, _, err := gateClient.GetTaskByID(context.Background(), executionID)
	if err != nil {
		fmt.Printf("failed to get deployment task details: %s\n", err.Error())
		os.Exit(1)
	}

	if execution.Status != "SUCCEEDED" {
		fmt.Printf("deployment failed.\n")
		os.Exit(1)
	}

	fmt.Println("deployment succeeded")
}

func providerFactory(manifest, manifestPath string) (ops.ManifestProvider, error) {
	if manifestPath == "" && manifest == "" {
		return nil, fmt.Errorf("must provide manifest or manifestPath")
	}

	if manifestPath != "" {
		fmt.Println("manifest path provided")
		return ops.ManifestsFromPath(manifestPath)
	}
	return ops.ManifestsFromString(manifestPath)
}

func performDeploy(task ops.Task, client gate.Client) (string, error) {
	output, _, err := client.CreateTask(context.Background(), task)
	if err != nil {
		return "", err
	}
	parts := strings.Split(output.Ref, "/")
	executionID := parts[len(parts)-1]
	return executionID, nil
}

func pollExecution(ctx context.Context, executionID string, client gate.Client) error {
	var topErr error
	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			execution, _, err := client.GetTaskByID(ctx, executionID)
			if err != nil {
				topErr = err
				break
			}
			if isCompleteStatus(execution.Status) {
				return nil
			}
		}
	}

	return topErr
}

var (
	statusSuccess  = "SUCCEEDED"
	statusTerminal = "TERMINAL"
)

func isCompleteStatus(status string) bool {
	return status == statusSuccess || status == statusTerminal
}
