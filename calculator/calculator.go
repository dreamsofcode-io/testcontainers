package calculator

import (
	"context"
	"database/sql"
	"fmt"
)

type Calculator struct {
	db *sql.DB
}

func New(db *sql.DB) *Calculator {
	return &Calculator{
		db: db,
	}
}

func (c *Calculator) Add(ctx context.Context, a int, b int) (int, error) {
	var res int
	err := c.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %d + %d", a, b)).Scan(&res)
	return res, err
}
