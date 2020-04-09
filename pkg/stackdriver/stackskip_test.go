// +build unit

package stackdriver

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

type logWrapper struct {
	Logger *logrus.Logger
}

func (l *logWrapper) error(msg string) {
	l.Logger.Error(msg)
}

func TestStackSkip(t *testing.T) {
	var out bytes.Buffer

	logger := logrus.New()
	logger.Out = &out
	logger.Formatter = NewFormatter(
		WithService("test"),
		WithVersion("0.1"),
		WithStackSkip("github.com/jenkins-x/logrus-stackdriver-formatter/pkg/stackdriver"),
	)

	mylog := logWrapper{
		Logger: logger,
	}

	mylog.error("my log entry")

	var got map[string]interface{}
	json.Unmarshal(out.Bytes(), &got)
	got["timestamp"] = "2020-01-01T00:00:00.000000Z"

	want := map[string]interface{}{
		"severity":  "ERROR",
		"message":   "my log entry",
		"timestamp": "2020-01-01T00:00:00.000000Z",
		"serviceContext": map[string]interface{}{
			"service": "test",
			"version": "0.1",
		},
		"context": map[string]interface{}{
			"reportLocation": map[string]interface{}{
				"filePath":     "testing/testing.go",
				"lineNumber":   865.0,
				"functionName": "tRunner",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected output = %# v; want = %# v", pretty.Formatter(got), pretty.Formatter(want))
	}
}
