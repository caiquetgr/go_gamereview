package tests

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("\tStarting API integration tests")
	m.Run()
}
