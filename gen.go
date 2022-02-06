package promtail

//go:generate ./bin/mockgen --destination=httpClient/http_client_test.go --source=httpClient/client.go --package=httpClient

//go:generate ./bin/mockgen --destination=client/promtail_mocks_test.go --source=types.go --package=client
