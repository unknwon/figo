package core

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/Unknwon/com"
	"github.com/fsouza/go-dockerclient"

	"github.com/Unknwon/figo/modules/base"
	"github.com/Unknwon/figo/modules/log"
)

type (
	Options map[string]map[string]interface{}

	Link struct {
		*Service
		Name string
	}
	Links   map[string]Link
	Volumes map[string]interface{}

	Service struct {
		name    string
		client  *docker.Client
		project string
		links   Links
		volumes Volumes
		options map[string]interface{}
	}
)

func NewService(
	name string,
	client *docker.Client,
	project string,
	links Links,
	volumes Volumes,
	options map[string]interface{}) *Service {
	return &Service{
		name:    name,
		client:  client,
		project: project,
		links:   links,
		volumes: volumes,
		options: options,
	}
}

func (s *Service) GetLinkedNames() []string {
	links := make([]string, 0, len(s.links))
	for s := range s.links {
		links = append(links, s)
	}
	return links
}

// CanBeBuilt returns true if this is buildable service.
func (s *Service) CanBeBuilt() bool {
	_, ok := s.options["build"]
	return ok
}

// buildTagName returns the tag to give to images built for this service.
func (s *Service) buildTagName() string {
	return s.project + "_" + s.name
}

var imageIdPattern = regexp.MustCompile("Successfully built ([0-9a-f]+)")

func (s *Service) Build(noCache bool) (string, error) {
	log.Info("Building %s...", s.name)

	dockerfile := path.Join(s.options["build"].(string), "Dockerfile")
	if !com.IsFile(dockerfile) {
		return "", fmt.Errorf("build dockerfile does not exist or is not a file: %s", dockerfile)
	}

	file, err := os.Open(dockerfile)
	if err != nil {
		return "", fmt.Errorf("fail to open dockerfile: %v", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("fail to read dockerfile: %v", err)
	}
	fi, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("fail to get dockerfile info: %v", err)
	}
	inputbuf := bytes.NewBuffer(nil)
	so := base.NewStreamOutput()

	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: int64(len(data)), ModTime: fi.ModTime()})
	tr.Write(data)
	tr.Close()
	opts := docker.BuildImageOptions{
		Name:           s.buildTagName(),
		NoCache:        noCache,
		RmTmpContainer: true,
		InputStream:    inputbuf,
		OutputStream:   so,
		RawJSONStream:  true,
	}
	if err := s.client.BuildImage(opts); err != nil {
		return "", err
	}

	var imageId string
	if len(so.Events) > 0 {
		e := so.Events[len(so.Events)-1]
		m := imageIdPattern.FindAllStringSubmatch(e["stream"], 1)
		if m != nil {
			imageId = m[0][1]
		}
	}

	return imageId, nil
}

// HasApiContainer returns true if the container was created to fulfill this service.
func (s *Service) HasApiContainer(apiContainer *docker.APIContainers, oneOff bool) bool {
	name := base.GetApiContainerName(apiContainer)
	if len(name) == 0 || !base.IsValidApiContainerName(name, oneOff) {
		return false
	}
	projectName, serviceName, _ := base.ParseApiContainerName(name)
	return projectName == s.project && serviceName == s.name
}

// Containers returns a list of containers belong to service.
func (s *Service) Containers(stopped, oneOff bool) ([]*Container, error) {
	apiContainers, err := s.client.ListContainers(docker.ListContainersOptions{All: stopped})
	if err != nil {
		return nil, fmt.Errorf("fail to list containers: %v", err)
	}
	containers := make([]*Container, 0, len(apiContainers))
	for _, apiContainer := range apiContainers {
		if s.HasApiContainer(&apiContainer, oneOff) {
			containers = append(containers, NewContainerFromPs(s.client, &apiContainer))
		}
	}
	return containers, nil
}

// CreateContainer creates a container for this service.
// If the image doesn't exist, attempt to pull it.
func (s *Service) CreateContainer(oneOff bool) (*Container, error) {
	return nil, nil
}

// StartContainer starts a existing container.
func (s *Service) StartContainer(c, intermediate *Container) error {
	if c == nil {

	}
	return nil
}

func (s *Service) StartContainerIfStopped(c *Container) error {
	if c.IsRunning() {
		return nil
	}
	log.Info("Starting %s...", c.Name)
	return s.StartContainer(c, nil)
}

func (s *Service) Start() error {
	containers, err := s.Containers(true, false)
	if err != nil {
		return fmt.Errorf("fail to get containers(%s): %v", s.name, err)
	}

	// TODO
	for _, c := range containers {
		if err = s.StartContainerIfStopped(c); err != nil {
			return fmt.Errorf("fail to start container(%s): %v", c.Name, err)
		}
	}
	return nil
}
