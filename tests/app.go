package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//nolint:unused
type appContainer struct {
	testcontainers.Container
	URI string
}

//nolint:unused
func startAppContainer(ctx context.Context, networkName string) (*appContainer, error) {
	dockerFile := ".docker/app/Dockerfile"
	contextPath := "../"

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Dockerfile: dockerFile,
			Context:    contextPath,
			KeepImage:  true,
		},
		ExposedPorts: []string{"8081/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(10 * time.Second).WithPort("8081/tcp"),
		Networks: []string{
			networkName,
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	var appC *appContainer

	if container != nil {
		appC = &appContainer{Container: container}
	}
	if err != nil {
		return appC, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return appC, err
	}

	mappedPort, err := container.MappedPort(ctx, "8081/tcp")
	if err != nil {
		return appC, err
	}

	appC.URI = fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	return appC, nil
}
