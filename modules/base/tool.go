package base

import (
	"os"
)

// DockerUrl returns Docker API URL from environment variable.
func DockerUrl() string {
	return os.Getenv("DOCKER_HOST")
}
