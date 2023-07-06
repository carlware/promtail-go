package client

import (
	"time"

	"github.com/carlware/promtail-go"
)

type optionApplyFunc func(client *promClient) error

type Option interface {
	applyOption(client *promClient) error
}

func (f optionApplyFunc) applyOption(p *promClient) error {
	return f(p)
}

func WithStreamConverter(converter promtail.StreamConverter) Option {
	return optionApplyFunc(func(client *promClient) error {
		client.streamConv = converter
		return nil
	})
}

func WithStaticLabels(labels map[string]interface{}) Option {
	return optionApplyFunc(func(client *promClient) error {
		client.staticLabels = labels
		return nil
	})
}

func WithWriteTimeout(ms int) Option {
	return optionApplyFunc(func(client *promClient) error {
		client.writeTimeout = time.Duration(ms) * time.Millisecond
		return nil
	})
}
