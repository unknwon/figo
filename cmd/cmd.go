package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"

	"github.com/Unknwon/figo/modules/base"
	"github.com/Unknwon/figo/modules/core"
	"github.com/Unknwon/figo/modules/log"
)

// GetProjectName returns given project name,
// if it is empty, then guesses it by config file path,
// otherwise, just returns 'defualt'.
func GetProjectName(configPath, name string) string {
	// TODO: normalize_name
	if len(name) > 0 {
		return name
	}
	absPath, _ := filepath.Abs(configPath)
	name = path.Base(path.Dir(absPath))
	if len(name) > 0 {
		return name
	}
	return "default"
}

// GetConfig loads and returns fig configuration.
func GetConfig(configPath string) (core.Options, error) {
	if !com.IsExist(configPath) {
		return nil, core.FigFileNotFound{configPath}
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := make(core.Options)
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func GetClient(verbose bool) (*docker.Client, error) {
	baseUrl := base.DockerUrl()
	client, err := docker.NewClient(baseUrl)
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Info("Figo version: %s", base.AppVer)
		log.Info("Docker host: %s", baseUrl)
		// env, err := client.Version()
		// if err != nil {
		// 	log.Fatal("Fail to get docker version: %v", err)
		// }
		// log.Info("Docker version: %s", env)
	}
	return client, nil
}

func GetProject(name, configPath string, verbose bool) (*core.Project, error) {
	config, err := GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config file: %v", err)
	}
	client, err := GetClient(verbose)
	if err != nil {
		return nil, fmt.Errorf("fail to create new client: %v", err)
	}
	return core.NewProjectFromConfig(GetProjectName(configPath, name), config, client)
}

func Setup(ctx *cli.Context) (*core.Project, error) {
	log.Verbose = ctx.GlobalBool("verbose")
	return GetProject(ctx.GlobalString("project-name"), ctx.GlobalString("file"), log.Verbose)
}
