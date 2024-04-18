package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/testcontainers/database"
	"github.com/dreamsofcode-io/testcontainers/repository"
)

func TestCreate(t *testing.T) {
	checkParallel(t)

	testCases := []struct {
		name   string
		setup  func(ctx context.Context, db *sql.DB) error
		input  repository.CreateData
		errors bool
	}{
		{
			name: "happy path",
			input: repository.CreateData{
				Name:   "firebolt",
				Damage: 100,
				Mana:   10,
			},
			errors: false,
		},
		{
			name: "empty name",
			input: repository.CreateData{
				Name:   "",
				Damage: 200,
				Mana:   20,
			},
			errors: true,
		},
		{
			name: "name collision",
			setup: func(ctx context.Context, conn *sql.DB) error {
				repo := repository.New(conn)
				_, err := repo.Create(ctx, repository.CreateData{
					Name:   "icewheel",
					Damage: 100,
					Mana:   15,
				})
				return err
			},
			input: repository.CreateData{
				Name:   "icewheel",
				Damage: 200,
				Mana:   20,
			},
			errors: true,
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
				tc.setup(ctx, conn)
			}

			now := time.Now().Truncate(time.Millisecond)
			spell, err := repo.Create(ctx, tc.input)
			time.Sleep(sleepTime)
			after := time.Now().Truncate(time.Millisecond)

			if tc.errors {
				assert.Error(t, err)

				// Ensure nothing exists in database
				row := conn.QueryRowContext(
					ctx,
					"SELECT id FROM spell WHERE name = $1 AND damage = $2 AND mana = $3",
					tc.input.Name,
					tc.input.Damage,
					tc.input.Mana,
				)

				var id string
				err := row.Scan(&id)

				assert.ErrorIs(t, err, sql.ErrNoRows)

				return
			}
			// Assert no error
			assert.NoError(t, err)

			// Check spell properties
			assert.Equal(t, spell.Name, tc.input.Name)
			assert.Equal(t, spell.Mana, tc.input.Mana)
			assert.Equal(t, spell.Damage, tc.input.Damage)
			assert.Equal(t, spell.CreatedAt, spell.UpdatedAt)
			assert.GreaterOrEqual(t, spell.CreatedAt, now)
			assert.LessOrEqual(t, spell.CreatedAt, after)

			// Ensure row exists in database
			row := conn.QueryRowContext(
				ctx,
				"SELECT id, name, mana, damage, created_at, updated_at FROM spell WHERE id = $1",
				spell.ID,
			)

			var rowSpell repository.Spell
			err = row.Scan(
				&rowSpell.ID,
				&rowSpell.Name,
				&rowSpell.Mana,
				&rowSpell.Damage,
				&rowSpell.CreatedAt,
				&rowSpell.UpdatedAt,
			)

			assert.NoError(t, err)
			assert.Equal(t, rowSpell, spell)
		})
	}
}
