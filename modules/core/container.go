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

func (c *Container) Name() string {
	return c.dict["Name"]
}

// Return a value from the container or None if the value is not set.
//:param key: a string using dotted notation for nested dictionary lookups
func (c *Container) Get(key string) {

}

func (c *Container) IsRunning() bool {
	return false
}
