package filemanager

import (
	"bytes"
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"
)

type FileManager struct {
	ctx context.Context
}

func NewFileManager(ctx context.Context) *FileManager {
	return &FileManager{
		ctx: ctx,
	}
}

func (f *FileManager) newHTTPRequest(url string, r *http.Request) *http.Request {
	prxReq, _ := http.NewRequestWithContext(f.ctx, r.Method, url, r.Body)
	prxQuery := prxReq.URL.Query()

	for key, values := range r.URL.Query() {
		for _, value := range values {
			prxQuery.Add(key, value)
		}
	}

	for key, values := range r.Header {
		for _, value := range values {
			prxReq.Header.Set(key, value)
		}
	}

	prxReq.URL.RawQuery = prxQuery.Encode()

	return prxReq
}

func (f *FileManager) FetchFile(url string, r *http.Request) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}

	req := f.newHTTPRequest(url, r)

	slog.Debug(fmt.Sprintf("Proxy: IN='%s %s' -> OUT='%s %s'", r.Method, r.URL.String(), req.Method, req.URL.String())) //nolint

	res, err := http.DefaultClient.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("error fetching remote http image: %w", err)
	}

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			slog.Error(fmt.Sprintf("unable to close response body: %s", closeErr))
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(
			fmt.Sprintf(
				"error fetching remote http image: (status=%d) (url=%s)",
				res.StatusCode,
				req.URL.String(),
			),
			res.StatusCode,
		)
	}

	// Read the body
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to create image from response body: %w (url=%s)", err, req.URL.String())
	}

	return buf, nil
}

func (f *FileManager) SaveFile(path string, body []byte) error {
	_ = f
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			slog.Error(fmt.Sprintf("unable to close os.File: %s", closeErr))
		}
	}(file)

	reader := bytes.NewReader(body)

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileManager) ReadFile(filePath string) ([]byte, error) {
	_ = f
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error read file: %w", err)
	}

	return content, nil
}

func (f *FileManager) DeleteFile(filePath string) error {
	_ = f
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}

func (f *FileManager) GetFileNameByURL(url string) string {
	_ = f
	ext := path.Ext(url)
	hash := md5.Sum([]byte(url)) //nolint:gosec
	return fmt.Sprintf("%s%s", hex.EncodeToString(hash[:]), ext)
}
