package rest_service

import (
	"bytes"
	"net/http"
	"time"
)

type RestInterface interface {
	Get(data RestQueryParams) (error, *http.Response)
	Put(data RestQueryParams) (error, *http.Response)
	Post(data RestQueryParams) (error, *http.Response)
}

type RestQueryParams struct {
	Ep      string
	Headers *http.Header
	Body    *bytes.Buffer
}

type RestResponse struct {
	Body []byte
}

type RestInstanceInitParams struct {
	Domain  string
	Timeout time.Duration
}
