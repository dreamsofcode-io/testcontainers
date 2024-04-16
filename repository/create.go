package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const createQuery string = `
INSERT INTO spell (id, name, damage, mana, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
`

type CreateData struct {
	Name   string
	Damage int
	Mana   uint
}

func (r *Spells) Create(ctx context.Context, data CreateData) (Spell, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Spell{}, fmt.Errorf("failed to generate uuid v7")
	}

	timestamp := time.Now()

	_, err = r.db.ExecContext(
		ctx,
		createQuery,
		id,
		data.Name,
		data.Damage,
		data.Mana,
		timestamp,
	)
	if err != nil {
		return Spell{}, fmt.Errorf("failed to exec context: %w", err)
	}

	return Spell{
		ID:        id,
		Name:      data.Name,
		Damage:    data.Damage,
		Mana:      data.Mana,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}, nil
}
