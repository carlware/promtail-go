package promtail

import "context"

type Entry []interface{}

// LabelSet is a key/value pair mapping of labels
type LabelSet map[string]interface{}

// Stream represents a log stream.  It includes a set of log entries and their labels.
type Stream struct {
	Labels  LabelSet `json:"stream"`
	Entries []Entry  `json:"values"`
}

// PushRequest models a log stream push
type PushRequest struct {
	Streams []*Stream `json:"streams"`
}

type HttpClient interface {
	Push(ctx context.Context, request PushRequest) error
}

type StreamConverter interface {
	ConvertEntry([]byte) (Entry, error)
	ExtractLabels([]byte) (LabelSet, error)
}
