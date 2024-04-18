package calculator_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/dreamsofcode-io/testcontainers/calculator"
	"github.com/dreamsofcode-io/testcontainers/database"
)

type input struct {
	a int
	b int
}

func TestAdd(t *testing.T) {
	ctx := context.Background()

	request := testcontainers.ContainerRequest{
		Image: "postgres:16",
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		t.Fatal("failed to start container:", err)
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatal("failed to get endpoint:", err)
	}

	connURI := fmt.Sprintf("postgres://user:secret@%s/testdb?sslmode=disable", endpoint)

	testCases := []struct {
		name  string
		input input
		wants int
	}{
		{
			name: "Simple 1 + 1",
			input: input{
				a: 1,
				b: 1,
			},
			wants: 2,
		},
		{
			name: "Simple 5 + 6",
			input: input{
				a: 5,
				b: 6,
			},
			wants: 11,
		},
		{
			name: "Simple 10 - 12",
			input: input{
				a: 10,
				b: -12,
			},
			wants: -2,
		},
	}

	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {

			conn, err := database.Connect(connURI)
			assert.NoError(t, err)

			calc := calculator.New(conn)
			res, err := calc.Add(ctx, tc.input.a, tc.input.b)

			assert.NoError(t, err)
			assert.Equal(t, tc.wants, res)
		})
	}
}
