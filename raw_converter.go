package promtail

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// taken from https://github.com/acarl005/stripansi/blob/master/stripansi.go
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func strip(str string) string {
	return re.ReplaceAllString(str, "")
}

type rawStreamConv struct {
	labels map[string]bool
	sep    string
}

func NewRawStreamConv(labels, sep string) StreamConverter {
	labelsArr := strings.Split(labels, ",")
	labelsMap := make(map[string]bool, len(labelsArr))
	for _, label := range labelsArr {
		labelsMap[label] = true
	}

	return &rawStreamConv{
		labels: labelsMap,
		sep:    sep,
	}
}

func (s *rawStreamConv) ConvertEntry(bytes []byte) (Entry, error) {
	now := time.Now().UnixNano()
	return Entry{strconv.FormatInt(now, 10), string(bytes)}, nil
}

func (s *rawStreamConv) ExtractLabels(bytes []byte) (LabelSet, error) {
	if len(s.labels) == 0 {
		return LabelSet{}, nil
	}
	if s.sep == "" {
		return LabelSet{}, nil
	}

	labelSet := LabelSet{}
	rawText := strip(string(bytes))

	tokens := strings.Split(rawText, " ")
	for _, token := range tokens {
		if strings.Contains(token, s.sep) {
			kv := strings.Split(token, s.sep)
			if len(kv) == 2 {
				key := kv[0]
				value := kv[1]
				if s.labels[key] {
					labelSet[key] = value
				}
			}
		}
	}
	return labelSet, nil
}
