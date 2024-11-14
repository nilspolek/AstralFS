package functionservice

import (
	"context"
	"io"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type svc struct {
	apiClient client.APIClient
	functions fns
}

type fns []Function

func New() (FunctionService, error) {
	var (
		out *svc
		err error
	)
	out = &svc{}
	out.apiClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (fs *svc) CreateFunction(fn Function) (int, error) {
	var port nat.Port = nat.Port(strconv.Itoa(fn.Port) + "/tcp")

	containerCfg := &container.Config{
		Image: fn.Image,
		ExposedPorts: nat.PortSet{
			port: {},
		},
	}
	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			port: {},
		},
	}

	// Pull Image
	readCloser, err := fs.apiClient.ImagePull(context.TODO(), fn.Image, image.PullOptions{})
	if err != nil {
		return 0, err
	}
	defer readCloser.Close()
	_, err = io.Copy(io.Discard, readCloser)

	// Create Container
	containerResp, err := fs.apiClient.ContainerCreate(context.Background(), containerCfg, hostCfg, nil, nil, fn.Id.String())
	if err != nil {
		return 0, err
	}
	fs.functions = append(fs.functions, fn)

	// Start Container
	err = fs.apiClient.ContainerStart(context.TODO(), containerResp.ID, container.StartOptions{})
	if err != nil {
		return 0, err
	}

	// Get the container's port bindings (host port)
	containerJSON, err := fs.apiClient.ContainerInspect(context.Background(), containerResp.ID)
	if err != nil {
		return 0, err
	}

	// Extract the host port from the container's PortBindings
	hostPort := ""
	for _, bindings := range containerJSON.NetworkSettings.Ports {
		if len(bindings) > 0 {
			hostPort = bindings[0].HostPort
			break
		}
	}
	return strconv.Atoi(hostPort)
}

func (fs *svc) DeleteFunction(functionId uuid.UUID) error {
	if err := fs.apiClient.ContainerStop(context.TODO(), functionId.String(), container.StopOptions{}); err != nil {
		return err
	}
	if err := fs.apiClient.ContainerRemove(context.TODO(), functionId.String(), container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return err
	}
	fs.functions = fs.functions.deleteByUUID(functionId)
	return nil
}

func (fs *svc) GetFunctions() ([]Function, error) {
	return fs.functions, nil
}

func (fs svc) Close() error {
	return fs.Close()
}

func (fns *fns) deleteByUUID(id uuid.UUID) fns {
	res := make([]Function, len(*fns)-1)
	j := 0
	for _, fn := range *fns {
		if fn.Id == id {
			continue
		}
		res[j] = fn
		j++
	}
	return res
}
