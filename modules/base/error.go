package base

import (
	"fmt"
)

type NoSuchService struct {
	Name string
}

func (e NoSuchService) Error() string {
	return fmt.Sprintf("no such service: %s", e.Name)
}

type FigFileNotFound struct {
	Name string
}

func (e FigFileNotFound) Error() string {
	return fmt.Sprintf("can't find %s. Are you in the right directory?", e.Name)
}

type ConfigurationError struct {
	Msg string
}

func (e ConfigurationError) Error() string {
	return e.Msg
}

type DependencyError struct {
	Msg string
}

func (e DependencyError) Error() string {
	return e.Msg
}
