package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UpdateData CreateData

const updateQuery string = `
UPDATE spell SET name = $1, damage = $2, mana = $3, updated_at = $4
WHERE id = $5
RETURNING	id, name, damage, mana, created_at, updated_at
`

func (r *Spells) Update(ctx context.Context, id uuid.UUID, updateData UpdateData) (Spell, error) {
	now := time.Now()

	res := r.db.QueryRowContext(
		ctx,
		updateQuery,
		updateData.Name,
		updateData.Damage,
		updateData.Mana,
		now,
		id,
	)

	spell := Spell{}

	err := res.Scan(
		&spell.ID,
		&spell.Name,
		&spell.Damage,
		&spell.Mana,
		&spell.CreatedAt,
		&spell.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Spell{}, ErrNotFound
	}
	if err != nil {
		return Spell{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return spell, nil
}
