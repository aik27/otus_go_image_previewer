package proxy

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
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

func (p *Client) FetchFile(url string, r *http.Request) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("https://%s", url)
	}

	req := p.newHTTPRequest(url, r)

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

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to create image from response body: %w (url=%s)", err, req.URL.String())
	}

	return buf, nil
}
