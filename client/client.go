package client

import (
	"context"
	"net/http"
	"time"

	"github.com/carlware/promtail-go"
	promtailHttpClient "github.com/carlware/promtail-go/httpClient"
)

const (
	readTimeOut        = 10
	writeTimeOut       = 10
	maxIdleConnections = 128
	maxConnections     = 512
)

func newHttpClient() *http.Client {
	customTransport := &(*http.DefaultTransport.(*http.Transport))

	customTransport.MaxConnsPerHost = maxConnections
	customTransport.MaxIdleConnsPerHost = maxIdleConnections
	customTransport.MaxIdleConns = maxIdleConnections

	return &http.Client{
		Transport: customTransport,
		Timeout:   readTimeOut * time.Second,
	}
}

type promClient struct {
	httpClient   promtail.HttpClient
	streamConv   promtail.StreamConverter
	staticLabels map[string]interface{}
	writeTimeout time.Duration
}

func NewSimpleClient(host, username, password string, opts ...Option) (*promClient, error) {
	pHttpClient, err := promtailHttpClient.New(host, username, password,
		promtailHttpClient.WithHttpClient(newHttpClient()),
	)
	if err != nil {
		return nil, err
	}
	client := &promClient{
		streamConv: promtail.NewRawStreamConv("", ""),
		httpClient: pHttpClient,
	}

	for _, opt := range opts {
		if err = opt.applyOption(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *promClient) Write(p []byte) (i int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.writeTimeout)
	defer cancel()

	labels, err := c.streamConv.ExtractLabels(p)
	if err != nil {
		return 0, err
	}

	entry, err := c.streamConv.ConvertEntry(p)
	if err != nil {
		return 0, err
	}

	for k, v := range c.staticLabels {
		labels[k] = v
	}

	req := promtail.PushRequest{Streams: []*promtail.Stream{{
		Labels:  labels,
		Entries: []promtail.Entry{entry},
	}}}

	if rErr := c.httpClient.Push(ctx, req); rErr != nil {
		return 0, rErr
	}

	return len(p), nil
}
