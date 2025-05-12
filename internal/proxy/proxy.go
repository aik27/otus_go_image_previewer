package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var (
	ErrHTTPRequestError              = errors.New("http request error")
	ErrHTTPRequestUnexpectedStatus   = errors.New("unexpected HTTP status")
	ErrUnableToCloseResponseBody     = errors.New("unable to close response body")
	ErrUnableToReadImageFromResponse = errors.New("unable to read image from response body")
)

type Client struct {
	ctx context.Context
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		ctx: ctx,
	}
}

func (p *Client) newHTTPRequest(url string, r *http.Request) *http.Request {
	prxReq, _ := http.NewRequestWithContext(p.ctx, r.Method, url, r.Body)
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

func (p *Client) FetchFile(url string, r *http.Request) ([]byte, int, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("http://%s", url)
	}

	req := p.newHTTPRequest(url, r)

	slog.Debug(fmt.Sprintf("Proxy: IN='%s %s' -> OUT='%s %s'", r.Method, r.URL.String(), req.Method, req.URL.String())) //nolint

	res, err := http.DefaultClient.Do(req) //nolint:bodyclose
	if err != nil {
		return nil,
			http.StatusInternalServerError,
			errors.Join(ErrHTTPRequestError, err)
	}

	defer func(Body io.ReadCloser) {
		if Body != nil {
			closeErr := Body.Close()
			if closeErr != nil {
				slog.Error(errors.Join(ErrUnableToCloseResponseBody, closeErr).Error())
			}
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil,
			res.StatusCode,
			fmt.Errorf(
				"%w: (status=%d) (url=%s)",
				ErrHTTPRequestUnexpectedStatus,
				res.StatusCode,
				req.URL.String(),
			)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil,
			http.StatusInternalServerError,
			errors.Join(ErrUnableToReadImageFromResponse, err)
	}

	return buf, http.StatusOK, nil
}
