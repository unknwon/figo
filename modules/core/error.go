package core

import (
	"fmt"
)

type NoSuchService struct {
	Name string
}

func (e NoSuchService) Error() string {
	return fmt.Sprintf("no such service: %s", e.Name)
}
