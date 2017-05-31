package cacheclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"time"
)

type Client struct {
	httpClient        *http.Client
	baseURL           string
	requestsPerSecond int
	rateLimiter       chan int
}

type apiConfig struct {
	host string
	path string
}

var defaultRequestsPerSecond = 10

func NewClient(host string) *Client {
	c := &Client{
		requestsPerSecond: defaultRequestsPerSecond,
		httpClient:        &http.Client{},
		baseURL:           host,
	}

	// Implement a bursty rate limiter.
	// Allow up to 1 second worth of requests to be made at once.
	c.rateLimiter = make(chan int, c.requestsPerSecond)
	// Prefill rateLimiter with 1 seconds worth of requests.
	for i := 0; i < c.requestsPerSecond; i++ {
		c.rateLimiter <- 1
	}
	go func() {
		// Wait a second for pre-filled quota to drain
		time.Sleep(time.Second)
		// Then, refill rateLimiter continuously
		for _ = range time.Tick(time.Second / time.Duration(c.requestsPerSecond)) {
			c.rateLimiter <- 1
		}
	}()

	return c
}

func (c *Client) get(ctx context.Context, config *apiConfig) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.rateLimiter:
		// Execute request.
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}
	req, err := http.NewRequest("GET", host+config.path, nil)
	if err != nil {
		return nil, err
	}

	return ctxhttp.Do(ctx, c.httpClient, req)
}

func (c *Client) post(ctx context.Context, config *apiConfig, body []byte) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.rateLimiter:
		// Execute request.
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}
	req, err := http.NewRequest("POST", host+config.path, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("post err ", err)
		return nil, err
	}

	return ctxhttp.Do(ctx, c.httpClient, req)
}

func (c *Client) delete(ctx context.Context, config *apiConfig) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.rateLimiter:
		// Execute request.
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}
	req, err := http.NewRequest("DELETE", host+config.path, nil)
	if err != nil {
		return nil, err
	}

	return ctxhttp.Do(ctx, c.httpClient, req)
}

func (c *Client) getJSON(ctx context.Context, config *apiConfig, resp interface{}) error {
	httpResp, err := c.get(ctx, config)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		respErr := map[string]string{}
		if err := json.NewDecoder(httpResp.Body).Decode(&respErr); err != nil {
			return err
		}
		return errors.New(respErr["error"])
	}

	return json.NewDecoder(httpResp.Body).Decode(&resp)
}

func (c *Client) postJSON(ctx context.Context, config *apiConfig, body []byte, resp interface{}) error {
	httpResp, err := c.post(ctx, config, body)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusCreated {
		respErr := map[string]string{}
		if err := json.NewDecoder(httpResp.Body).Decode(&respErr); err != nil {
			return err
		}

		return errors.New(respErr["error"])
	}

	return json.NewDecoder(httpResp.Body).Decode(&resp)
}

func (c *Client) deleteJSON(ctx context.Context, config *apiConfig, resp interface{}) error {
	httpResp, err := c.delete(ctx, config)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		respErr := map[string]string{}
		if err := json.NewDecoder(httpResp.Body).Decode(&respErr); err != nil {
			return err
		}
		return errors.New(respErr["error"])
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return err
	}

	return nil
}
