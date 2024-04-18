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

func TestUpdate(t *testing.T) {
	checkParallel(t)

	now := time.Now()

	setup := func(ctx context.Context, conn *sql.DB) error {
		_, err := conn.ExecContext(
			ctx,
			"INSERT INTO spell (id, name, mana, damage, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)",
			"c856b3a1-31fe-46ce-8823-40059d48a27c",
			"foo",
			10,
			40,
			now.Truncate(time.Millisecond),
		)
		if err != nil {
			return err
		}

		_, err = conn.ExecContext(
			ctx,
			"INSERT INTO spell (id, name, mana, damage, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)",
			"41fca60e-9b26-4196-9b18-d4a79b350523kd",
			"bar",
			70,
			800,
			now.Truncate(time.Millisecond),
		)
		return err
	}

	type input struct {
		data repository.UpdateData
		id   uuid.UUID
	}

	type want struct {
		spell repository.Spell
		err   error
	}

	testCases := []struct {
		name  string
		setup func(ctx context.Context, db *sql.DB) error
		input input
		wants want
	}{
		{
			name:  "happy path",
			setup: setup,
			input: input{
				id: uuid.MustParse("c856b3a1-31fe-46ce-8823-40059d48a27c"),
				data: repository.UpdateData{
					Name:   "firebolt",
					Damage: 100,
					Mana:   10,
				},
			},
			wants: want{
				spell: repository.Spell{
					ID:     uuid.MustParse("c856b3a1-31fe-46ce-8823-40059d48a27c"),
					Name:   "firebolt",
					Damage: 100,
					Mana:   10,
				},
			},
		},
		{
			name:  "missing spell",
			setup: setup,
			input: input{
				id: uuid.MustParse("0980dd52-bcc2-4019-9710-7816fc8c50bf"),
				data: repository.UpdateData{
					Name:   "firebolt",
					Damage: 100,
					Mana:   10,
				},
			},
			wants: want{
				err: repository.ErrNotFound,
			},
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

			time.Sleep(sleepTime)

			spell, err := repo.Update(ctx, tc.input.id, tc.input.data)

			if tc.wants.err != nil {
				assert.ErrorIs(t, err, tc.wants.err)
				return
			} else {
				// Assert no error
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wants.spell.Name, spell.Name)
			assert.Equal(t, tc.wants.spell.Damage, spell.Damage)
			assert.Equal(t, tc.wants.spell.Mana, spell.Mana)

			assert.Greater(t, spell.UpdatedAt, spell.CreatedAt)

			repoSpell, err := repo.FindByID(ctx, spell.ID)
			assert.NoError(t, err)

			assert.Equal(t, spell, repoSpell)
		})
	}
}
