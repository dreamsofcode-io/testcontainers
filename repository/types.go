package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Spells struct {
	db *sql.DB
}

func New(db *sql.DB) *Spells {
	return &Spells{
		db: db,
	}
}

type Spell struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Damage    int       `json:"damage"`
	Mana      uint      `json:"mana"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNotFound = errors.New("spell not found for id")
