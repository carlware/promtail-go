package promtail

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Strip(t *testing.T) {
	msg := "\x1b[38;5;140mfoo\x1b[0m bar"
	res := strip(msg)

	assert.Equal(t, res, "foo bar")
}

func Test_ExtractLabels(t *testing.T) {
	tests := []struct {
		name        string
		labels      string
		sep         string
		input       []byte
		expected    LabelSet
		expectedErr string
	}{
		{
			name:     "should return an empty label set if label string is blank",
			labels:   "",
			input:    []byte("info method=POST bytes_in=0 bytes_out=108"),
			expected: LabelSet{},
		},
		{
			name:     "should return an empty label set if sep is blank",
			labels:   "method",
			input:    []byte("info method=POST bytes_in=0 bytes_out=108"),
			expected: LabelSet{},
		},
		{
			name:   "should return a label set with only label tags",
			labels: "method,bytes_in",
			sep:    "=",
			input:  []byte("info method=POST bytes_in=0 bytes_out=108"),
			expected: LabelSet{
				"method":   "POST",
				"bytes_in": "0",
			},
		},
		{
			name:   "should return a label set with only label tags and strip all terminal colors",
			labels: "method,bytes_in",
			sep:    "=",
			input:  []byte("\u001B[38;5;140minfo\u001B[0m method=POST bytes_in=0 bytes_out=108"),
			expected: LabelSet{
				"method":   "POST",
				"bytes_in": "0",
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rsc := NewRawStreamConv(tc.labels, tc.sep)
			got, err := rsc.ExtractLabels(tc.input)
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
