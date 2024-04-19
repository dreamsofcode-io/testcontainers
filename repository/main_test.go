package repository_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/dreamsofcode-io/testcontainers/database"
)

var connURL = ""
var parrallel = false
var sleepTime = time.Millisecond * 500

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("foobar"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalln("failed to load container:", err)
	}

	connURL, err = container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalln("failed to get connection string:", err)
	}

	migrate, err := database.Migrate(connURL)
	if err != nil {
		log.Fatal("failed to migrate db: ", err)
	}

	res := m.Run()

	migrate.Drop()

	os.Exit(res)
}

func getConnection(ctx context.Context) (*sql.DB, error) {
	return database.Connect(connURL)
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
