package tests

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("\tStarting consumer integration tests")
	m.Run()
}
