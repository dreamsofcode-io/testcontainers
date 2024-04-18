package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dreamsofcode-io/spellbook/database"
	"github.com/dreamsofcode-io/spellbook/repository"
)

func TestCreate(t *testing.T) {
	testCases := []struct{}{}
	conn, err := database.Connect(
		"postgresql://user:secret@localhost:5432/testdb?sslmode=disable",
	)
	assert.NoError(t, err)

	repo := repository.New(conn)

	ctx := context.Background()
	spell, err := repo.Create(ctx, tc.input)

	if tc.wants.err != nil {
		assert.ErrorIs(t, err, tc.wants.err)
	} else {
		assert.NoError(t, err)
	}

	assert.Equal(t, spell, tc.wants.spell)
}
