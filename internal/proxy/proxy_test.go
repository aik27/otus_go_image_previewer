package proxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchFile_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test content"))
	}))
	defer server.Close()

	client := NewClient(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	data, status, err := client.FetchFile(server.URL, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, []byte("test content"), data)
}

func TestFetchFile_RequestError(t *testing.T) {
	t.Parallel()

	client := NewClient(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, status, err := client.FetchFile("http://invalid-url", req)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, status)
	require.True(t, errors.Is(err, ErrHTTPRequestError))
}

func TestFetchFile_HTTPStatusError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, status, err := client.FetchFile(server.URL, req)
	require.Error(t, err)
	require.Equal(t, http.StatusNotFound, status)
	require.True(t, errors.Is(err, ErrHTTPRequestUnexpectedStatus))
}
