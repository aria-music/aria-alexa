package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aria-music/aria-alexa/aria"
)

var client = newRestClient("")

type restClient struct {
	client http.Client
	token  string
}

func newRestClient(token string) *restClient {
	if token == "" {
		token = ariaToken
	}

	return &restClient{
		token: token,
		client: http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (r *restClient) doAriaRequest(ctx context.Context, request *aria.Request) (*aria.Response, error) {
	enc, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ariaEndpoint, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect endpoint: %w", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status error: %d", resp.StatusCode)
	}

	aresp, err := aria.ParseResponse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return aresp, err
}

func (r *restClient) sendOP(ctx context.Context, op string) error {
	req := &aria.Request{
		Token: r.token,
		OP:    op,
	}

	_, err := r.doAriaRequest(ctx, req)
	return err
}

func (r *restClient) sendRequest(ctx context.Context, op string, data interface{}) (*aria.Response, error) {
	req := &aria.Request{
		Token: r.token,
		OP:    op,
		Data:  data,
	}

	return r.doAriaRequest(ctx, req)
}

// dispatch dispatches
func sendOP(ctx context.Context, op string) error {
	return client.sendOP(ctx, op)
}

func sendRequest(ctx context.Context, op string, data interface{}) (*aria.Response, error) {
	return client.sendRequest(ctx, op, data)
}
