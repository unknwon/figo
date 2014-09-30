package core

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"

	"github.com/Unknwon/figo/modules/log"
)

// Project represents a collection of services.
type Project struct {
	name     string
	services []*Service
	client   *docker.Client
}

// NewProject initializes and returns a minimal project.
func NewProject(name string, services []*Service, client *docker.Client) *Project {
	return &Project{
		name:     name,
		services: services,
		client:   client,
	}
}

func NewProjectFromDicts(name string, dicts []map[string]interface{}, client *docker.Client) (*Project, error) {
	pro := NewProject(name, []*Service{}, client)
	// FIXME: loop sort_service_dicts
	// for _,
	return pro, nil
}

func NewProjectFromConfig(name string, config Options, client *docker.Client) (*Project, error) {
	dicts := make([]map[string]interface{}, 0, len(config))
	for name, service := range config {
		if service == nil {
			return nil, ConfigurationError{fmt.Sprintf("Service \"%s\" doesn't have any configuration options. All top level keys in your fig.yml must map to a dictionary of configuration options", name)}
		}
		service["name"] = name
		dicts = append(dicts, service)
	}
	return NewProjectFromDicts(name, dicts, client)
}

// ListServicesNames returns a list of services' names.
func (p *Project) ListServicesNames() []string {
	names := make([]string, len(p.services))
	for i, s := range p.services {
		names[i] = s.name
	}
	return names
}

// GetService retrieve a service by name.
// It returns NoSuchService if the named service does not exist.
func (p *Project) GetService(name string) (*Service, error) {
	for _, s := range p.services {
		if s.name == name {
			return s, nil
		}
	}
	return nil, NoSuchService{name}
}

func (p *Project) injectLinks(services []*Service) (_ []*Service, err error) {
	allServices := make([]*Service, 0, len(services))
	for _, s := range services {
		var linkedServices []*Service
		linkedNames := s.GetLinkedNames()
		if len(linkedNames) > 0 {
			linkedServices, err = p.Services(linkedNames, true)
			if err != nil {
				return nil, err
			}
			allServices = append(allServices, linkedServices...)
		}
	}
	return allServices, nil
}

// Services returns a list of this project's services filtered
// by the provided list of entries, or all services if entries is empty or nil.
//
// If includeLinks is true, returns a list including the links for
// entries, in order of dependency.
//
// Preserves the original order of Project.services where possible,
// reordering as needed to resolve links.
//
// It returns NoSuchService if any of the named services do not exist.
func (p *Project) Services(entries []string, includeLinks bool) (_ []*Service, err error) {
	// Return all services.
	if entries == nil || len(entries) == 0 {
		return p.Services(p.ListServicesNames(), includeLinks)
	}

	services := make([]*Service, len(entries))
	for i, name := range entries {
		s, err := p.GetService(name)
		if err != nil {
			return nil, err
		}
		services[i] = s
	}

	if includeLinks {
		services, err = p.injectLinks(services)
		if err != nil {
			return nil, err
		}
	}

	cache := map[string]bool{}
	uniques := make([]*Service, 0, len(services))
	for _, s := range services {
		if !cache[s.name] {
			uniques = append(uniques, s)
			cache[s.name] = true
		}
	}
	return uniques, nil
}

func (p *Project) Build(entries []string, noCache bool) error {
	services, err := p.Services(entries, false)
	if err != nil {
		return err
	}
	for _, s := range services {
		if s.CanBeBuilt() {
			s.Build(noCache)
		} else {
			log.Info("%s uses an image, skipping", s.name)
		}
	}
	return nil
}
