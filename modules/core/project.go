package core

import (
	"fmt"
	"strings"

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

func in(slice []interface{}, str string) bool {
	for _, s := range slice {
		if str == strings.Split(s.(string), ":")[0] {
			return true
		}
	}
	return false
}

func visit(dicts, unmarked Options, marked map[string]bool, sorted *[]map[string]interface{}, dict map[string]interface{}) error {
	name := dict["name"].(string)
	if marked[name] {
		if links, ok := dict["links"]; ok && in(links.([]interface{}), name) {
			return DependencyError{fmt.Sprintf("A service can not link to itself: %s", name)}
		}
		if links, ok := dict["volumes_from"]; ok && in(links.([]interface{}), name) {
			return DependencyError{fmt.Sprintf("A service can not mount itself as volume: %s", name)}
		}
		return DependencyError{fmt.Sprintf("Circular import found: %s", name)}
	}
	if _, ok := unmarked[name]; ok {
		marked[name] = true
		if links, ok := dict["links"]; ok {
			for _, link := range links.([]interface{}) {
				dict := dicts[strings.Split(link.(string), ":")[0]]
				if dict == nil {
					return DependencyError{fmt.Sprintf("Link service does not exist: %s", link)}
				}
				if err := visit(dicts, unmarked, marked, sorted, dict); err != nil {
					return err
				}
			}
		}
		if links, ok := dict["volumes_from"]; ok {
			for _, link := range links.([]interface{}) {
				dict := dicts[strings.Split(link.(string), ":")[0]]
				if dict == nil {
					return DependencyError{fmt.Sprintf("Link service does not exist: %s", link)}
				}
				if err := visit(dicts, unmarked, marked, sorted, dict); err != nil {
					return err
				}
			}
		}
		delete(marked, name)
		delete(unmarked, name)
		*sorted = append((*sorted)[:1], (*sorted)[:]...)
		(*sorted)[0] = dict
	}
	return nil
}

func SortServiceDicts(dicts Options) ([]map[string]interface{}, error) {
	unmarked := make(Options)
	for k, v := range dicts {
		unmarked[k] = v
	}
	marked := make(map[string]bool)
	sorted := make([]map[string]interface{}, 0, len(dicts))

	for _, dict := range dicts {
		if err := visit(dicts, unmarked, marked, &sorted, dict); err != nil {
			return nil, err
		}
	}
	return sorted, nil
}

func (p *Project) GetLinks(dict map[string]interface{}) map[string]string {
	links := make(map[string]string)
	var (
		linkStr     string
		serviceName string
		linkName    string
	)
	if rawLinks, ok := dict["links"]; ok {
		for _, link := range rawLinks.([]interface{}) {
			linkStr = link.(string)
			if strings.Contains(linkStr, ":") {
				infos := strings.SplitN(linkStr, ":", 2)
				serviceName = infos[0]
				linkName = infos[1]
			} else {
				serviceName = linkStr
			}
			// FIXME: get_service
			links[serviceName] = linkName
		}
		delete(dict, "links")
	}
	return links
}

func (p *Project) GetVolumesFrom(dict map[string]interface{}) map[string]string {
	volumes := make(map[string]string)
	if rawVolumes, ok := dict["volumes_from"]; ok {
		for _, volume := range rawVolumes.([]interface{}) {
			volumeName := volume.(string)
			// FIXME: get_service
			// FIXME: Container.from_id(self.client, volume_name)
			volumes[volumeName] = ""
		}
		delete(dict, "volumes_from")
	}
	return volumes
}

func NewProjectFromDicts(name string, dicts Options, client *docker.Client) (*Project, error) {
	pro := NewProject(name, []*Service{}, client)
	sorted, err := SortServiceDicts(dicts)
	if err != nil {
		return nil, err
	}
	for _, dict := range sorted {
		serviceName := dict["name"].(string)
		links := pro.GetLinks(dict)
		volumes := pro.GetVolumesFrom(dict)
		pro.services = append(pro.services, NewService(serviceName, client, name, links, volumes, dict))
	}
	return pro, nil
}

func NewProjectFromConfig(name string, config Options, client *docker.Client) (*Project, error) {
	dicts := make(Options)
	for name, service := range config {
		if service == nil {
			return nil, ConfigurationError{fmt.Sprintf("Service \"%s\" doesn't have any configuration options. All top level keys in your fig.yml must map to a dictionary of configuration options", name)}
		}
		service["name"] = name
		dicts[name] = service
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
			linkedServices, err = p.GetServices(linkedNames, true)
			if err != nil {
				return nil, err
			}
			allServices = append(allServices, linkedServices...)
		}
	}
	return allServices, nil
}

// GetServices returns a list of this project's services filtered
// by the provided list of entries, or all services if entries is empty or nil.
//
// If includeLinks is true, returns a list including the links for
// entries, in order of dependency.
//
// Preserves the original order of Project.services where possible,
// reordering as needed to resolve links.
//
// It returns NoSuchService if any of the named services do not exist.
func (p *Project) GetServices(entries []string, includeLinks bool) (_ []*Service, err error) {
	// Return all services.
	if entries == nil || len(entries) == 0 {
		return p.GetServices(p.ListServicesNames(), includeLinks)
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

	set := map[string]bool{}
	uniques := make([]*Service, 0, len(services))
	for _, s := range services {
		if !set[s.name] {
			uniques = append(uniques, s)
			set[s.name] = true
		}
	}
	return uniques, nil
}

func (p *Project) Build(entries []string, noCache bool) error {
	services, err := p.GetServices(entries, false)
	if err != nil {
		return err
	}
	for _, s := range services {
		if s.CanBeBuilt() {
			if _, err = s.Build(noCache); err != nil {
				return err
			}
		} else {
			log.Info("%s uses an image, skipping", s.name)
		}
	}
	return nil
}
