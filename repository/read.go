package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const baseFindQuery = `
SELECT id, name, damage, mana, created_at, updated_at
FROM spell
`

var findByIDQuery = fmt.Sprintf(
	"%s WHERE id = $1", baseFindQuery,
)

func (r *Spells) FindByID(ctx context.Context, id uuid.UUID) (Spell, error) {
	row := r.db.QueryRowContext(ctx, findByIDQuery, id)

	res := Spell{}
	err := row.Scan(&res.ID, &res.Name, &res.Damage, &res.Mana, &res.CreatedAt, &res.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Spell{}, ErrNotFound
	} else if err != nil {
		return Spell{}, ErrNotFound
	}

	return res, nil
}

func (r *Spells) FindAll(ctx context.Context) ([]Spell, error) {
	var spells []Spell

	rows, err := r.db.QueryContext(ctx, baseFindQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		spell := Spell{}
		err = rows.Scan(
			&spell.ID,
			&spell.Name,
			&spell.Damage,
			&spell.Mana,
			&spell.CreatedAt,
			&spell.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		spells = append(spells, spell)
	}

	return spells, nil
}
