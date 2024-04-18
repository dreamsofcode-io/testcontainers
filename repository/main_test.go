package repository_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/dreamsofcode-io/testcontainers/database"
)

var connURL = "postgresql://user:secret@localhost:5432/testdb?sslmode=disable"
var parrallel = false
var sleepTime = time.Millisecond * 500

func TestMain(m *testing.M) {
	migrate, err := database.Migrate(connURL)
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}

	res := m.Run()

	migrate.Drop()

	os.Exit(res)
}

func cleanup() {
	conn, err := database.Connect(connURL)
	if err != nil {
		return
	}

	conn.Exec("DELETE FROM spell")
}

func checkParallel(t *testing.T) {
	if parrallel {
		t.Parallel()
	}
}
