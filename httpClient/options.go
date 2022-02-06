package httpClient

type optionApplyFunc func(client *promHttpClient) error

type Option interface {
	applyOption(client *promHttpClient) error
}

func (f optionApplyFunc) applyOption(p *promHttpClient) error {
	return f(p)
}

func WithHttpClient(c HttpDoer) Option {
	return optionApplyFunc(func(client *promHttpClient) error {
		client.httpClient = c
		return nil
	})
}
