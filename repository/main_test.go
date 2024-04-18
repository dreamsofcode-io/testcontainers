package repository_test

import (
	"log"
	"os"
	"testing"

	"github.com/dreamsofcode-io/testcontainers/database"
)

var connURL = "postgresql://user:secret@localhost:5432/testdb?sslmode=disable"

func TestMain(m *testing.M) {
	migrate, err := database.Migrate(connURL)
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}

	res := m.Run()

	migrate.Drop()

	os.Exit(res)
}
