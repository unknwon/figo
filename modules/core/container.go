package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Unknwon/com"
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

func valOrNil(val interface{}) string {
	if val == nil {
		return ""
	}
	return val.(string)
}

func parseBool(str string) bool {
	val, _ := strconv.ParseBool(str)
	return val
}

// CreateContainer creates new container by given options.
func CreateContainer(client *docker.Client, options map[string]interface{}) (*Container, error) {
	// FIXME: many things need to be done.
	// https://docs.docker.com/reference/commandline/cli/#run
	createOptions := docker.CreateContainerOptions{
		Name: options["name"].(string),
		Config: &docker.Config{
			Hostname:     valOrNil(options["hostname"]),
			Domainname:   valOrNil(options["domainname"]),
			User:         valOrNil(options["user"]),
			Memory:       com.StrTo(valOrNil(options["memory"])).MustInt64(), // FIXME: parse unit?
			CpuShares:    com.StrTo(valOrNil(options["cpu-shares"])).MustInt64(),
			CpuSet:       valOrNil(options["cpuset"]),
			AttachStdin:  strings.Contains(valOrNil(options["attach"]), "STDIN"),
			AttachStdout: strings.Contains(valOrNil(options["attach"]), "STDOUT"),
			AttachStderr: strings.Contains(valOrNil(options["attach"]), "STDERR"),
			Tty:          parseBool(valOrNil(options["tty"])),
			OpenStdin:    parseBool(valOrNil(options["interactive"])),
			// Env:             valOrNil(options["env"]),
			// Dns:             valOrNil(options["dns"]),
			VolumesFrom: valOrNil(options["volumes-from"]),
			WorkingDir:  valOrNil(options["workdir"]),
			// Entrypoint:      valOrNil(options["entrypoint"]),
		},
	}

	// TODO: PortSpecs, ExposedPorts, Volumes
	c, err := client.CreateContainer(createOptions)
	if err != nil {
		return nil, err
	}
	return NewContainerFromId(client, c.ID)
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
	case "NetworkSettings.Ports":
		return c.NetworkSettings.Ports
	}
	return nil
}

// IsRunning returns true if container is running.
func (c *Container) IsRunning() bool {
	running, ok := c.Get("State.Running").(bool)
	return ok && running
}

func (c *Container) Ports() map[docker.Port][]docker.PortBinding {
	ports, ok := c.Get("NetworkSettings.Ports").(map[docker.Port][]docker.PortBinding)
	if !ok {
		return map[docker.Port][]docker.PortBinding{}
	}
	return ports
}

func (c *Container) Stop() error {
	return c.client.StopContainer(c.ID, 60)
}

func (c *Container) Start() error {
	return c.client.StartContainer(c.ID, &docker.HostConfig{})
}

func (c *Container) Wait() (int, error) {
	return c.client.WaitContainer(c.ID)
}

func (c *Container) Kill() error {
	return c.client.KillContainer(docker.KillContainerOptions{})
}

func (c *Container) GetLocalPort(port int, protocol string) string {
	portBinding := c.Ports()[docker.Port(fmt.Sprintf("%d/%s", port, protocol))]
	return fmt.Sprintf("%s:%s", portBinding[0].HostIp, portBinding[0].HostPort)
}

func (c *Container) Restart() error {
	return c.client.RestartContainer(c.ID, 60)
}
