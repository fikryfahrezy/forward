package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments SET
			deleted_at = NOW()
		WHERE
			id = $1
			AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
