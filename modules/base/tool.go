package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// DockerUrl returns Docker API URL from environment variable.
func DockerUrl() string {
	url := os.Getenv("DOCKER_HOST")
	if len(url) > 0 {
		return url
	}
	return "unix:///var/run/docker.sock"
}

type StreamOutput struct {
	Events []map[string]string
}

func NewStreamOutput() *StreamOutput {
	return &StreamOutput{
		Events: make([]map[string]string, 0),
	}
}

func (so *StreamOutput) Write(p []byte) (int, error) {
	e := make(map[string]string)
	for _, line := range bytes.Split(p, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		if err := json.Unmarshal(line, &e); err != nil {
			return 0, err
		}
		fmt.Print(string(e["stream"]))
		so.Events = append(so.Events, e)
	}
	return len(p), nil
}
