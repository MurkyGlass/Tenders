package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]models.Role, error)
	GetByID(ctx context.Context, id int) (*models.Role, error)
}

type roleRepository struct {
	db QueryExecutor
}

func NewRoleRepository(db QueryExecutor) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAll(ctx context.Context) ([]models.Role, error) {
	var role []models.Role
	query := `SELECT * FROM role`
	err := r.db.SelectContext(ctx, &role, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return role, nil
}
func (r *roleRepository) GetByID(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role

	query := `SELECT * FROM role WHERE id_role = $1`

	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}
