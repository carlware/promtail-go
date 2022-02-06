package httpClient

import (
	"context"
	"fmt"
	"github.com/carlware/promtail-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func getMockRespPayload(body string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(body))
}

func TestPromTailClient_Push(t *testing.T) {
	const host = "loki.com"
	const pass = "secret"
	const user = "loki"

	type mockSetters struct {
		httpClient func(mock *MockHttpDoer)
	}

	tests := []struct {
		name        string
		req         promtail.PushRequest
		mockSetters mockSetters
		expectedErr string
	}{
		{
			name: "should return with a proper request",
			req: promtail.PushRequest{Streams: []*promtail.Stream{
				{
					Labels: promtail.LabelSet{
						"foo": "bar",
					},
					Entries: []promtail.Entry{
						{
							1,
							"info method=POST",
						},
					},
				},
			}},
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpDoer) {
					mock.EXPECT().
						Do(gomock.Any()).
						Return(&http.Response{
							StatusCode: http.StatusOK,
							Body:       getMockRespPayload(``),
						}, nil)
				},
			},
		},
		{
			name: "should return an error if httpClient fails",
			req:  promtail.PushRequest{},
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpDoer) {
					mock.EXPECT().
						Do(gomock.Any()).
						Return(nil, fmt.Errorf("http client fails"))
				},
			},
			expectedErr: "http client fails",
		},
		{
			name: "should return an error if server response return an error code",
			req:  promtail.PushRequest{},
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpDoer) {
					mock.EXPECT().
						Do(gomock.Any()).
						Return(&http.Response{
							StatusCode: http.StatusUnauthorized,
							Body:       getMockRespPayload(`unauthorized`),
						}, nil)
				},
			},
			expectedErr: "error: unauthorized",
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			httpClientMock := NewMockHttpDoer(ctrl)
			tc.mockSetters.httpClient(httpClientMock)

			client, err := New(host, user, pass, WithHttpClient(httpClientMock))
			require.NoError(t, err)

			err = client.Push(context.Background(), tc.req)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
