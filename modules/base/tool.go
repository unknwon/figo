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

	"github.com/Unknwon/figo/modules/log"
)

var (
	DockerConfigKeys = []string{"image", "command", "hostname", "domainname", "user", "detach", "stdin_open", "tty", "mem_limit", "ports", "environment", "dns", "volumes", "entrypoint", "privileged", "volumes_from", "net", "working_dir"}
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

// IsValidContainerName returns true if given name is a valid container name.
func IsValidContainerName(name string, oneOff bool) bool {
	m := apiContainerNamePattern.FindAllStringSubmatch(name, -1)
	if m == nil {
		return false
	}
	if oneOff {
		return m[0][2] == "run_"
	}
	return len(m[0][2]) == 0
}

// ParseContainerName parses and returns container name.
func ParseContainerName(name string) (string, string, int) {
	m := apiContainerNamePattern.FindAllStringSubmatch(name, -1)
	return m[0][0], m[0][1], com.StrTo(m[0][3]).MustInt()
}

// ParseArgs returns options and services' name from command line arguments.
// FIXME: parse slice values
func ParseArgs(args []string) (entries []string, _ map[string]string) {
	log.Warn("ParseArgs has limitations and bugs! Cannot handle all the arguments.")
	options := map[string]string{}
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			infos := strings.SplitN(arg, "=", 2)
			options[strings.TrimLeft(infos[0], "-")] = infos[1]
		} else {
			entries = append(entries, arg)
		}
	}
	return entries, options
}

// ParseVolumeSpec parses and returns given volume configuration.
func ParseVolumeSpec(config string) ([]string, error) {
	infos := strings.Split(config, ":")
	if len(infos) > 3 {
		return nil, ConfigurationError{fmt.Sprintf("Volume %s has incorrect format, should be external:internal[:mode]", config)}
	} else if len(infos) == 1 {
		return []string{"", infos[0], "rw"}, nil
	}

	if len(infos) == 2 {
		infos = append(infos, "rw")
	}

	if infos[2] != "rw" && infos[2] != "ro" {
		return nil, ConfigurationError{fmt.Sprintf("Volume %s has invalid mode (%s), should be one of: rw, ro", config, infos[2])}
	}
	return infos, nil
}
