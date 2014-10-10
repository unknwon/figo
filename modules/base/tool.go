package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Unknwon/com"
	"github.com/fsouza/go-dockerclient"
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

// GetApiContainerName returns name of API container.
func GetApiContainerName(apiContainer *docker.APIContainers) string {
	for _, name := range apiContainer.Names {
		infos := strings.Split(name, "/")
		if len(infos) == 2 {
			return infos[1]
		}
	}
	return ""
}

var apiContainerNamePattern = regexp.MustCompile(`^([^_]+)_([^_]+)_(run_)?(\d+)$`)

// IsValidApiContainerName returns true if given name is a valid API container name.
func IsValidApiContainerName(name string, oneOff bool) bool {
	m := apiContainerNamePattern.FindAllStringSubmatch(name, -1)
	if m == nil {
		return false
	}
	if oneOff {
		return m[0][2] == "run_"
	}
	return len(m[0][2]) == 0
}

// ParseApiContainerName parses and returns API container name.
func ParseApiContainerName(name string) (string, string, int) {
	m := apiContainerNamePattern.FindAllStringSubmatch(name, -1)
	return m[0][0], m[0][1], com.StrTo(m[0][3]).MustInt()
}
