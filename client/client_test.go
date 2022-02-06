package client

import (
	"fmt"
	"github.com/carlware/promtail-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Client_Write(t *testing.T) {
	const host = "loki.com"
	const pass = "secret"
	const user = "loki"

	type mockSetters struct {
		httpClient func(mock *MockHttpClient)
		streamConverter func(mock *MockStreamConverter)
	}

	tests := []struct{
		name string
		input []byte
		mockSetters mockSetters
		expected int
		expectedErr string
	}{
		{
			name: "should return no errors with a successfully request",
			input: []byte(`level=info method=POST bytes_in=10`),
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpClient) {
					mock.EXPECT().
						Push(gomock.Any(), gomock.Any()).
						Return(nil)
				},
				streamConverter: func(mock *MockStreamConverter) {
					mock.EXPECT().
						ExtractLabels(gomock.Any()).
						Return(promtail.LabelSet{}, nil)

					mock.EXPECT().
						ConvertEntry(gomock.Any()).
						Return(promtail.Entry{}, nil)
				},
			},
			expected: 34,
		},
		{
			name: "should return an error if stream converter ExtractLabels returns an error",
			input: []byte(`level=info method=POST bytes_in=10`),
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpClient) {

				},
				streamConverter: func(mock *MockStreamConverter) {
					mock.EXPECT().
						ExtractLabels(gomock.Any()).
						Return(promtail.LabelSet{}, fmt.Errorf("some error"))

				},
			},
			expectedErr: "some error",
		},
		{
			name: "should return an error if stream converter ConvertEntry returns an error",
			input: []byte(`level=info method=POST bytes_in=10`),
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpClient) {

				},
				streamConverter: func(mock *MockStreamConverter) {
					mock.EXPECT().
						ExtractLabels(gomock.Any()).
						Return(promtail.LabelSet{}, nil)

					mock.EXPECT().
						ConvertEntry(gomock.Any()).
						Return(promtail.Entry{}, fmt.Errorf("some error"))
				},
			},
			expectedErr: "some error",
		},
		{
			name: "should return an error if promtail http client return an error",
			input: []byte(`level=info method=POST bytes_in=10`),
			mockSetters: mockSetters{
				httpClient: func(mock *MockHttpClient) {
					mock.EXPECT().
						Push(gomock.Any(), gomock.Any()).
						Return(fmt.Errorf("some error"))
				},
				streamConverter: func(mock *MockStreamConverter) {
					mock.EXPECT().
						ExtractLabels(gomock.Any()).
						Return(promtail.LabelSet{}, nil)

					mock.EXPECT().
						ConvertEntry(gomock.Any()).
						Return(promtail.Entry{}, nil)
				},
			},
			expectedErr: "some error",
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := NewSimpleClient(host, user, pass)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)

			promtailHttpClientMock := NewMockHttpClient(ctrl)
			streamConverterMock := NewMockStreamConverter(ctrl)

			c.httpClient = promtailHttpClientMock
			c.streamConv = streamConverterMock

			tc.mockSetters.streamConverter(streamConverterMock)
			tc.mockSetters.httpClient(promtailHttpClientMock)

			got, err := c.Write(tc.input)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)

		})

	}
}
