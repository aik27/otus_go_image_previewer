//go:build integration

package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestIntegrationContainers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	newNetwork, err := network.New(ctx)
	require.NoError(t, err)
	testcontainers.CleanupNetwork(t, newNetwork)

	// -------------------------------------------
	// Up Nginx
	// -------------------------------------------

	nginxC, err := startNginxContainer(ctx, newNetwork.Name)
	testcontainers.CleanupContainer(t, nginxC)
	require.NoError(t, err, "failed to start nginx container")

	resp, err := http.Get(nginxC.URI) //nolint:all
	require.NoError(t, err, "failed HTTP GET to nginx container")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			require.NoError(t, err)
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected HTTP status code of nginx container")

	// -------------------------------------------
	// Up App
	// -------------------------------------------

	appC, appErr := startAppContainer(ctx, newNetwork.Name)
	testcontainers.CleanupContainer(t, appC)
	require.NoError(t, appErr, "failed to start app container")

	resp, appErr = http.Get(appC.URI) //nolint:all
	require.NoError(t, appErr, "failed HTTP GET to app container")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			require.NoError(t, err)
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected HTTP status code of app container")

	// -------------------------------------------
	// Table tests
	// -------------------------------------------

	startNoCached := time.Now()
	t.Run("Get images", func(t *testing.T) {
		imgGoAround(t, appC)
	})
	elapsedNoCached := time.Since(startNoCached)

	startCached := time.Now()
	t.Run("Get cached images", func(t *testing.T) {
		imgGoAround(t, appC)
	})
	elapsedCached := time.Since(startCached)

	require.Less(t, float64(elapsedCached), float64(elapsedNoCached/5), "some images are not cached")
}

func imgGoAround(t *testing.T, appC *appContainer) {
	t.Run("Go around", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			url := fmt.Sprintf("%s/fill/300/200/nginx:80/photo/0%d.jpg", appC.URI, i)

			resp, err := http.Get(url) //nolint:all
			require.NoError(t, err, "failed to get image from app container")
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					require.NoError(t, err)
				}
			}(resp.Body)

			require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected HTTP status code of app container")

			actualImage, _ := io.ReadAll(resp.Body)
			expectedImage, _ := os.ReadFile(fmt.Sprintf("testdata/0%d.jpg", i))

			require.Equal(t, expectedImage, actualImage, "unexpected image content")
		}
	})
}
