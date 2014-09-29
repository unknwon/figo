package core

import (
	"github.com/fsouza/go-dockerclient"

	"github.com/Unknwon/figo/modules/log"
)

type Service struct {
	name    string
	client  *docker.Client
	project string
	links   map[string]string
	options map[string]interface{}
}

func NewService(
	name string,
	links map[string]string,
	client *docker.Client,
	project string,
	options map[string]interface{}) *Service {
	return &Service{
		name:    name,
		client:  client,
		project: project,
		links:   links,
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

func (s *Service) CanBeBuilt() bool {
	_, ok := s.options["build"]
	return ok
}

// buildTagName returns the tag to give to images built for this service.
func (s *Service) buildTagName() string {
	return s.project + "_" + s.name
}

func (s *Service) Build(noCache bool) (string, error) {
	log.Info("Building %s...", s.name)

	// FIXME: Remote, InputStream, OutputStream
	opts := docker.BuildImageOptions{
		Name:           s.buildTagName(),
		NoCache:        noCache,
		RmTmpContainer: true,
		InputStream:    nil,
		OutputStream:   nil,
	}
	if err := s.client.BuildImage(opts); err != nil {
		return "", err
	}

	return "kao", nil
}
