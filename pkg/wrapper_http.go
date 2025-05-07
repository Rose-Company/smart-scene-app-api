package pkg

import (
	"encoding/json"
	"smart-scene-app-api/common"
	"time"

	"github.com/go-resty/resty/v2"
)

type WrapperConfig struct {
	retryTimes int
	timeout    time.Duration
}

type WrapperClient struct {
	*resty.Client
}

func newBaseHTTPClient(config *WrapperConfig) *resty.Client {
	return resty.New().
		SetDebug(true).
		SetTimeout(config.timeout).
		SetRetryCount(config.retryTimes).
		SetRetryWaitTime(1 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return err != nil && response.StatusCode() == 408 || response.StatusCode() == 404 || response.StatusCode() > 500
		})
}

func NewWrapperClient(config *WrapperConfig) *WrapperClient {
	if config == nil {
		config = &WrapperConfig{
			retryTimes: 0,
			timeout:    20 * time.Second,
		}
	}
	return &WrapperClient{Client: newBaseHTTPClient(config)}
}

func (c *WrapperClient) Get(url string, result interface{}) error {
	resp, err := c.Client.R().SetResult(&result).Get(url)
	if err != nil {
		return common.ErrCodeSystemErr
	}
	return c.handleResponse(resp, result)
}

func (c *WrapperClient) PostJSON(url string, body interface{}, result interface{}) error {
	resp, err := c.Client.R().SetBody(body).SetResult(&result).Post(url)
	if err != nil {
		return common.ErrCodeSystemErr
	}
	return c.handleResponse(resp, result)
}

func (c *WrapperClient) SetJWTAuth(token string) *WrapperClient {
	c.Client.SetAuthScheme("Bearer")
	c.Client.SetAuthToken(token)
	return c
}

func (c *WrapperClient) handleResponse(resp *resty.Response, result interface{}) error {
	if resp.IsError() {
		return common.ErrCodeSystemErr
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return common.ErrCodeSystemErr
	}
	return nil
}
