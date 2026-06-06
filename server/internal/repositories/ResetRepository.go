package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type ResetRepository interface {
	GetByToken(ctx context.Context, token string) (*models.ResetToken, error)
	Create(ctx context.Context, r *models.ResetToken) error
	DeleteByToken(ctx context.Context, token string) error
}

type resetRepository struct {
	db QueryExecutor
}

func NewResetRepository(db QueryExecutor) ResetRepository {
	return &resetRepository{db: db}
}

func (r *resetRepository) GetByToken(ctx context.Context, token string) (*models.ResetToken, error) {
	var t models.ResetToken

	query := `SELECT * FROM reset_tokens 
		WHERE token = $1 AND expires_at > NOW()`

	err := r.db.GetContext(ctx, &t, query, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return &t, nil
}
func (r *resetRepository) Create(ctx context.Context, refr *models.ResetToken) error {
	query := `
		INSERT INTO reset_tokens (token, id_user, expires_at)
		VALUES (:token, :id_user, :expires_at)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRowxContext(ctx, refr).Scan(&refr.ID)
	if err != nil {
		return fmt.Errorf("failed to create reset and get its ID: %w", err)
	}

	return nil
}
func (r *resetRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM reset_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}
