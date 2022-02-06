package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlware/promtail-go"
	"io/ioutil"
	"net/http"
)

const (
	readTimeOut        = 10
	maxIdleConnections = 128
	maxConnections     = 512
	path               = "/loki/api/v1/push"
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type promHttpClient struct {
	username   string
	password   string
	hostUrl    string
	httpClient HttpDoer
}

func New(host, username, password string, opts ...Option) (*promHttpClient, error) {
	pc := &promHttpClient{
		hostUrl:  host,
		username: username,
		password: password,
	}

	for _, opt := range opts {
		err := opt.applyOption(pc)
		if err != nil {
			return nil, err
		}
	}

	return pc, nil
}

func (p *promHttpClient) Push(ctx context.Context, req promtail.PushRequest) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, p.hostUrl+path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	request.Header.Add("content-type", "application/json")
	if p.password != "" {
		request.SetBasicAuth(p.username, p.password)
	}

	resp, err := p.httpClient.Do(request)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error: %s", string(respBody))
	}

	return nil
}
