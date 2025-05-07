package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"smart-scene-app-api/config"
	l "smart-scene-app-api/pkg/logger"
	"time"
)

const (
	DefaultHTTPContentType = "application/json"
	DefaultRequestTimeout  = 60
	DefaultRetryTimes      = 3
)

type HTTPClient interface {
	Post(url string, contentType string, headers map[string]string, data interface{}, retry int, timeout int) ([]byte, int, error)
}

type httpClient struct {
	client *http.Client
}

var (
	ll = l.New()
)

func NewHTTPClient() HTTPClient {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = config.Config.Http.MaxIdleConnection
	t.IdleConnTimeout = time.Duration(config.Config.Http.IdleConnectionTimeout) * time.Second
	client := &http.Client{
		Transport: t,
	}
	return &httpClient{
		client: client,
	}
}

func makeContext(timeout int) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = DefaultRequestTimeout
	}
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

func (h *httpClient) do(req *http.Request) ([]byte, int, error) {
	var (
		resp *http.Response
		errR error
	)

	resp, errR = h.client.Do(req)

	if errR != nil || resp.StatusCode != http.StatusOK {
		status := http.StatusInternalServerError
		if resp != nil {
			status = resp.StatusCode
		}
		return nil, status, errR
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, resp.StatusCode, err
}

func (h *httpClient) Post(
	url, contentType string,
	headers map[string]string,
	data interface{},
	retry, timeout int,
) ([]byte, int, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	ctx, cancel := makeContext(timeout)
	defer cancel()

	var (
		body   []byte
		status int
	)
	for i := 0; i < retry; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(dataByte))
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		if contentType == "" {
			contentType = DefaultHTTPContentType
		}
		req.Header.Set("Content-Type", contentType)
		if len(headers) > 0 {
			for header, value := range headers {
				req.Header.Add(header, value)
			}
		}

		body, status, err = h.do(req)
		if err == nil && status == http.StatusOK {
			break
		}
		if err != nil {
			ll.Error("Send request failed", l.Int("retry", i), l.Error(err))
		}
		time.Sleep(time.Duration(i*2) * time.Second)
	}
	return body, status, err
}
