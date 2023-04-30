package main

import (
	"fmt"
	"testing"

	"github.com/caiquetgr/go_gamereview/foundation/docker"
)

func Test_get_all_games(t *testing.T) {
}

func StartDB() (*docker.Container, error) {
	dbEnv := []string{"-e", "POSTGRES_DB=gamereview", "-e", "POSTGRES_USER=postgres", "-e", "POSTGRES_PASSWORD=postgres"}
	container, err := docker.StartContainer("postgres:15.2-alpine", "5432:5432", dbEnv...)
	if err != nil {
		return nil, fmt.Errorf("error running container: %v", err)
	}
	fmt.Println("started db container")
	return container, nil
}

func StopDB(db *docker.Container) error {
	return docker.StopContainer(db.ID)
}
