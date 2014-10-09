package core

import (
	"github.com/fsouza/go-dockerclient"
)

// Container represents a Docker container, constructed from the output of
// GET /containers/:id:/json.
type Container struct {
	client    *docker.Client
	dict      map[string]string
	inspected bool
}

// NewContainerFromId returns a container by given ID.
func NewContainerFromId(client *docker.Client, id string) (*docker.Container, error) {
	return client.InspectContainer(id)
}
