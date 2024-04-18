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

func TestFindByID(t *testing.T) {
	checkParallel(t)

	type want struct {
		err   error
		spell repository.Spell
	}

	now := time.Now()

	testCases := []struct {
		name  string
		setup func(ctx context.Context, conn *sql.DB) error
		input uuid.UUID
		wants want
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
			wants: want{
				err: nil,
				spell: repository.Spell{
					ID:        uuid.MustParse("c856b3a1-31fe-46ce-8823-40059d48a27c"),
					Name:      "foo",
					Mana:      10,
					Damage:    40,
					CreatedAt: now.Truncate(time.Millisecond),
					UpdatedAt: now.Truncate(time.Millisecond),
				},
			},
		},
		{
			name:  "non existant",
			input: uuid.MustParse("cea24ef2-c52c-45ed-a848-a8512012a830"),
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
				assert.NoError(t, tc.setup(ctx, conn))
			}

			spell, err := repo.FindByID(ctx, tc.input)
			time.Sleep(sleepTime)

			if tc.wants.err != nil {
				assert.ErrorIs(t, err, tc.wants.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, spell, tc.wants.spell)
		})
	}
}

func TestFindAll(t *testing.T) {
	checkParallel(t)

	now := time.Now().Truncate(time.Millisecond)

	testCases := []struct {
		name  string
		setup func(ctx context.Context, conn *sql.DB) error
		wants []repository.Spell
	}{
		{
			name: "empty repository",
		},
		{
			name: "some spells repository",
			setup: func(ctx context.Context, conn *sql.DB) error {
				spells := []repository.Spell{
					{
						ID:        uuid.MustParse("f3a88af4-bfb0-4981-b2c1-da45752148c9"),
						Name:      "firebolt",
						Mana:      10,
						Damage:    200,
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:        uuid.MustParse("1a738d7d-1b7d-429d-927b-f16547570625"),
						Name:      "magmalake",
						Mana:      90,
						Damage:    500,
						CreatedAt: now.Add(-time.Hour),
						UpdatedAt: now,
					},
				}

				for _, spell := range spells {
					_, err := conn.ExecContext(
						ctx,
						"INSERT INTO spell (id, name, mana, damage, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
						spell.ID,
						spell.Name,
						spell.Mana,
						spell.Damage,
						spell.CreatedAt,
						spell.UpdatedAt,
					)

					if err != nil {
						return err
					}
				}

				return nil
			},
			wants: []repository.Spell{
				{
					ID:        uuid.MustParse("f3a88af4-bfb0-4981-b2c1-da45752148c9"),
					Name:      "firebolt",
					Mana:      10,
					Damage:    200,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        uuid.MustParse("1a738d7d-1b7d-429d-927b-f16547570625"),
					Name:      "magmalake",
					Mana:      90,
					Damage:    500,
					CreatedAt: now.Add(-time.Hour),
					UpdatedAt: now,
				},
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
				assert.NoError(t, tc.setup(ctx, conn))
			}

			time.Sleep(sleepTime)

			spells, err := repo.FindAll(ctx)
			assert.NoError(t, err)

			assert.Equal(t, spells, tc.wants)
		})
	}
}
