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

func (s *Service) Containers(stopped, oneOff bool) ([]*Container, error) {
	return nil, nil
}

func (s *Service) StartContainer(c, intermediate *Container) error {
	return nil
}

func (s *Service) StartContainerIfStopped(c *Container) error {
	if c.IsRunning() {
		return nil
	}
	log.Info("Starting %s...", c.Name())
	return s.StartContainer(c, nil)
}

func (s *Service) Start() error {
	containers, err := s.Containers(true, false)
	if err != nil {
		return fmt.Errorf("fail to get containers: %v", err)
	}

	for _, c := range containers {
		if err = s.StartContainerIfStopped(c); err != nil {
			return fmt.Errorf("fail to start container(%s): %v", c.Name(), err)
		}
	}
	return nil
}
