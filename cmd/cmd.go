package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"

	"github.com/Unknwon/figo/modules/base"
	"github.com/Unknwon/figo/modules/core"
	"github.com/Unknwon/figo/modules/log"
)

// GetConfig loads and returns project configuration.
func GetConfig(cfgPath string) (core.Options, error) {
	if !com.IsExist(cfgPath) {
		return nil, base.FigFileNotFound{cfgPath}
	}
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	config := make(core.Options)
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// GetClient returns a new Docker client.
func GetClient(verbose bool) (*docker.Client, error) {
	baseUrl := base.DockerUrl()
	client, err := docker.NewClient(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("fail to create new client: %v", err)
	}
	if verbose {
		log.Info("Figo version: %s", base.AppVer)
		log.Info("Docker host: %s", baseUrl)
		env, err := client.Version()
		if err != nil {
			log.Fatal("Fail to get docker version: %v", err)
		}
		envs := []string(*env)
		sort.Strings(envs)
		log.Info("Docker version:\n%s", strings.Join(envs, "\n"))
	}
	return client, nil
}

var normalCharPattern = regexp.MustCompile("[a-zA-Z0-9]+")

// GetProjectName returns given project name,
// if it is empty, then guesses it by config file path,
// otherwise, just returns 'defualt'.
func GetProjectName(cfgPath, name string) string {
	name = normalCharPattern.FindString(name)
	if len(name) > 0 {
		return name
	}
	absPath, _ := filepath.Abs(cfgPath)
	name = path.Base(path.Dir(absPath))
	if len(name) > 0 {
		return name
	}
	return "default"
}

// GetProject initializes and returns a new project.
func GetProject(name, cfgPath string, verbose bool) (*core.Project, error) {
	config, err := GetConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config file: %v", err)
	}
	client, err := GetClient(verbose)
	if err != nil {
		return nil, fmt.Errorf("fail to create new client: %v", err)
	}
	return core.NewProjectFromConfig(GetProjectName(cfgPath, name), config, client)
}

func setup(ctx *cli.Context) (*core.Project, error) {
	log.Verbose = ctx.Bool("verbose")
	return GetProject(ctx.GlobalString("project-name"), ctx.GlobalString("file"), log.Verbose)
}
