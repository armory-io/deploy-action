package gate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"

	"github.com/spinnaker/spin/cmd/gateclient"
)

type CreateTaskResponse struct {
	Ref string `json:"ref"`
}

type GetTaskResponse struct {
	Status string `json:"status"`
}

type Client interface {
	CreateTask(ctx context.Context, task interface{}) (*CreateTaskResponse, *http.Response, error)
	GetTaskByID(ctx context.Context, taskID string) (*GetTaskResponse, *http.Response, error)
}

type gateClient struct {
	client *gateclient.GatewayClient
}

func (g *gateClient) CreateTask(ctx context.Context, task interface{}) (*CreateTaskResponse, *http.Response, error) {
	createResponse, resp, err := g.client.TaskControllerApi.TaskUsingPOST1(ctx, task)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		return nil, resp, errors.New("failed to create task")
	}

	var taskResponse CreateTaskResponse
	if err := mapstructure.Decode(createResponse, &taskResponse); err != nil {
		return nil, resp, fmt.Errorf("cound not map response: %w", err)
	}

	return &taskResponse, resp, nil
}

func (g *gateClient) GetTaskByID(ctx context.Context, taskID string) (*GetTaskResponse, *http.Response, error) {
	createResponse, resp, err := g.client.TaskControllerApi.GetTaskUsingGET1(ctx, taskID)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode != 200 {
		return nil, resp, errors.New("failed to get task details")
	}

	var taskResponse GetTaskResponse
	if err := mapstructure.Decode(createResponse, &taskResponse); err != nil {
		return nil, resp, fmt.Errorf("cound not map response: %w", err)
	}

	return &taskResponse, resp, nil
}

func New(configLocation, baseURL string) (*gateClient, error) {
	client, err := gateclient.NewGateClient(&noopUI{}, baseURL, "", configLocation, true)
	if err != nil {
		return nil, err
	}
	return &gateClient{client: client}, nil
}

// noopUI fulfills the cli.UI interface that the gate
// client depends on. we don't really care about
// those messages so we just blackhole them.
type noopUI struct{}

func (n *noopUI) Success(message string) {
}

func (n *noopUI) JsonOutput(data interface{}) {
}

func (n *noopUI) Ask(s string) (string, error) {
	return "", nil
}

func (n *noopUI) AskSecret(s string) (string, error) {
	return n.Ask(s)
}

func (n *noopUI) Output(s string) {
}

func (n *noopUI) Info(s string) {
}

func (n *noopUI) Error(s string) {
}

func (n *noopUI) Warn(s string) {
}
