package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const deleteQuery = `
DELETE FROM spell WHERE id = $1
`

func (r *Spells) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
