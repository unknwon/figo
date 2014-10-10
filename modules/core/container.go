package core

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"

	"github.com/Unknwon/figo/modules/base"
	"github.com/Unknwon/figo/modules/log"
)

// Container represents a Docker container, constructed from the output of
// GET /containers/:id:/json.
type Container struct {
	client *docker.Client
	*docker.Container
	inspected bool
}

// NewContainerFromId returns a container by given ID.
func NewContainerFromId(client *docker.Client, id string) (*Container, error) {
	c, err := client.InspectContainer(id)
	if err != nil {
		return nil, fmt.Errorf("fail to inspect container(%s): %v", id, err)
	}
	return &Container{client, c, true}, nil
}

//NewContainerFromPs returns a container object from the output of GET /containers/json.
func NewContainerFromPs(client *docker.Client, apiContainer *docker.APIContainers) *Container {
	return &Container{
		client,
		&docker.Container{
			ID:    apiContainer.ID,
			Image: apiContainer.Image,
			Name:  base.GetApiContainerName(apiContainer),
		},
		false,
	}
}

// Inspect inspects container information.
func (c *Container) Inspect() (err error) {
	c.Container, err = c.client.InspectContainer(c.ID)
	if err != nil {
		return err
	}
	c.inspected = true
	return nil
}

func (c *Container) InspectIfNotInspected() error {
	if !c.inspected {
		return c.Inspect()
	}
	return nil
}

// Return a value from the container or None if the value is not set.
//:param key: a string using dotted notation for nested dictionary lookups
func (c *Container) Get(key string) interface{} {
	if err := c.InspectIfNotInspected(); err != nil {
		log.Error("Fail to inspect container(%s): %v", c.Name, err)
		return ""
	}
	switch key {
	case "State.Running":
		return c.State.Running
	}
	return nil
}

// IsRunning returns true if container is running.
func (c *Container) IsRunning() bool {
	running, ok := c.Get("State.Running").(bool)
	return ok && running
}
