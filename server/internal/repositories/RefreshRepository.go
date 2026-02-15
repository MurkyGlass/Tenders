package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type RefreshRepository interface {
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	Create(ctx context.Context, r *models.RefreshToken) error
	Update(ctx context.Context, r *models.RefreshToken) error
	DeleteByToken(ctx context.Context, token string) error
}

type refreshRepository struct {
	db QueryExecutor
}

func NewRefreshRepository(db QueryExecutor) RefreshRepository {
	return &refreshRepository{db: db}
}

func (r *refreshRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var t models.RefreshToken

	query := `SELECT * FROM refresh_tokens 
		WHERE token = $1 AND expires_at > NOW()`

	err := r.db.GetContext(ctx, &t, query, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	return &t, nil
}
func (r *refreshRepository) Create(ctx context.Context, refr *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (token, id_user, expires_at)
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
		return fmt.Errorf("failed to create refresh and get its ID: %w", err)
	}

	return nil
}

func (r *refreshRepository) Update(ctx context.Context, t *models.RefreshToken) error {
	query := `
		UPDATE refresh_tokens 
		SET token = :token, expires_at = :expires_at 
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, t)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}
	return nil
}

func (r *refreshRepository) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}
