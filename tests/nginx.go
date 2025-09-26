package tests

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

//nolint:unused
type nginxContainer struct {
	testcontainers.Container
	URI string
}

//nolint:unused
func startNginxContainer(ctx context.Context, networkName string) (*nginxContainer, error) {
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
		Hostname: "nginx",
		Networks: []string{
			networkName,
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
