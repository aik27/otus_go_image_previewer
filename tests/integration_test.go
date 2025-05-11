//go:build integration

package tests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type nginxContainer struct {
	testcontainers.Container
	URI string
}

func startContainer(ctx context.Context) (*nginxContainer, error) {
	wwwRoot, _ := filepath.Abs("../.docker/nginx/var/www/html")
	nginxConf, _ := filepath.Abs("../.docker/nginx/etc/nginx/nginx.conf")
	domainConf, _ := filepath.Abs("../.docker/nginx/etc/nginx/conf.d/default.conf")

	req := testcontainers.ContainerRequest{
		Image:        "nginx:1.28",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(10 * time.Second),
		Mounts: testcontainers.ContainerMounts{
			testcontainers.BindMount(wwwRoot, "/var/www/html"),                     //nolint:staticcheck
			testcontainers.BindMount(nginxConf, "/etc/nginx/nginx.conf"),           //nolint:staticcheck
			testcontainers.BindMount(domainConf, "/etc/nginx/conf.d/default.conf"), //nolint:staticcheck
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	var nginxC *nginxContainer
	if container != nil {
		nginxC = &nginxContainer{Container: container}
	}
	if err != nil {
		return nginxC, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nginxC, err
	}

	mappedPort, err := container.MappedPort(ctx, "80")
	if err != nil {
		return nginxC, err
	}

	nginxC.URI = fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	return nginxC, nil
}

func TestIntegrationNginxLatestReturn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	nginxC, err := startContainer(ctx)
	testcontainers.CleanupContainer(t, nginxC)
	require.NoError(t, err)

	resp, err := http.Get(nginxC.URI + "/photo/05.jpg") //nolint:all
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			require.NoError(t, err)
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
