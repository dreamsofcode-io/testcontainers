package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/testcontainers/database"
	"github.com/dreamsofcode-io/testcontainers/repository"
)

func TestDelete(t *testing.T) {
	now := time.Now()
	checkParallel(t)

	testCases := []struct {
		name  string
		setup func(ctx context.Context, conn *sql.DB) error
		input uuid.UUID
		wants error
	}{
		{
			name: "happy path",
			setup: func(ctx context.Context, conn *sql.DB) error {
				_, err := conn.ExecContext(
					ctx,
					"INSERT INTO spell (id, name, mana, damage, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)",
					"c856b3a1-31fe-46ce-8823-40059d48a27c",
					"foo",
					10,
					40,
					now.Truncate(time.Millisecond),
				)
				return err
			},
			input: uuid.MustParse("c856b3a1-31fe-46ce-8823-40059d48a27c"),
			wants: nil,
		},
		{
			name:  "non existant",
			input: uuid.MustParse("cea24ef2-c52c-45ed-a848-a8512012a830"),
			wants: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkParallel(t)
			conn, err := database.Connect(connURL)
			assert.NoError(t, err)

			repo := repository.New(conn)
			t.Cleanup(cleanup)

			ctx := context.Background()

			if tc.setup != nil {
				assert.NoError(t, tc.setup(ctx, conn))
			}

			err = repo.Delete(ctx, tc.input)

			time.Sleep(sleepTime)

			if tc.wants != nil {
				assert.ErrorIs(t, err, tc.wants)
			} else {
				assert.NoError(t, err)
				// ensure spell no longer exists
				_, err = repo.FindByID(ctx, tc.input)
				assert.ErrorIs(t, err, repository.ErrNotFound)
			}
		})
	}
}
