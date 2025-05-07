package rest_service

import (
	"fmt"
	"net/http"
	"smart-scene-app-api/common"
	"time"
)

type apiClient struct {
	conn   *http.Client
	domain string
}

func NewRestInstance(initData RestInstanceInitParams) *apiClient {
	conn := &http.Client{Timeout: initData.Timeout * time.Millisecond}
	return &apiClient{conn: conn, domain: initData.Domain}
}

func (t *apiClient) Get(data RestQueryParams) (error, *http.Response) {
	fmt.Println("t.domain+data.Ep: ", t.domain+data.Ep)
	req, err := http.NewRequest("GET", t.domain+data.Ep, nil)
	if err != nil {
		return common.ErrorWrapper("failed to init request", err), nil
	}

	resp, err := t.conn.Do(req)
	if err != nil {
		return common.ErrorWrapper("failed to send request", err), nil
	}
	return nil, resp
}

func (t *apiClient) Put(data RestQueryParams) (error, *http.Response) {
	req, err := http.NewRequest(http.MethodPut, t.domain+data.Ep, data.Body)
	if err != nil {
		return common.ErrorWrapper("failed to init request", err), nil
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := t.conn.Do(req)
	if err != nil {
		return common.ErrorWrapper("failed to send request", err), nil
	}
	return nil, resp
}

func (t *apiClient) Post(data RestQueryParams) (error, *http.Response) {
	req, err := http.NewRequest(http.MethodPost, t.domain+data.Ep, data.Body)
	if err != nil {
		return common.ErrorWrapper("failed to init request", err), nil
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := t.conn.Do(req)
	if err != nil {
		return common.ErrorWrapper("failed to send request", err), nil
	}
	return nil, resp
}
