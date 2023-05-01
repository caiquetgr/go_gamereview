package dbtest

import (
	"fmt"

	"github.com/caiquetgr/go_gamereview/foundation/docker"
)

func StartDB() (*docker.Container, error) {
	dbEnv := []string{"-e", "POSTGRES_DB=gamereview", "-e", "POSTGRES_USER=postgres", "-e", "POSTGRES_PASSWORD=postgres"}
	var err error

	c, err := docker.StartContainer("postgres:15.2-alpine", "5432:5432", dbEnv...)
	if err != nil {
		return nil, fmt.Errorf("error running container: %v", err)
	}

	fmt.Println("started db container")
	return c, nil
}

func StopDB(db *docker.Container) {
	docker.StopContainer(db.ID)
}
