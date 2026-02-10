package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type RoleInCompanyRepository interface {
	GetAll(ctx context.Context) ([]models.RoleInCompany, error)
	GetByID(ctx context.Context, id int) (*models.RoleInCompany, error)
	Create(ctx context.Context, role *models.RoleInCompany) error
	Update(ctx context.Context, role *models.RoleInCompany) error
	Delete(ctx context.Context, id int) error
}

type roleincompanyRepository struct {
	db QueryExecutor
}

func NewRoleInCompanyRepository(db QueryExecutor) RoleInCompanyRepository {
	return &roleincompanyRepository{db: db}
}

func (r *roleincompanyRepository) GetAll(ctx context.Context) ([]models.RoleInCompany, error) {
	var role []models.RoleInCompany
	query := `SELECT * FROM Role_in_Company`
	err := r.db.SelectContext(ctx, &role, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return role, nil
}

func (r *roleincompanyRepository) GetByID(ctx context.Context, id int) (*models.RoleInCompany, error) {
	var role models.RoleInCompany

	query := `SELECT * FROM Role_in_Company WHERE id_role = $1`

	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}
func (r *roleincompanyRepository) Create(ctx context.Context, role *models.RoleInCompany) error {
	err := role.Validate()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO Role_in_Company (name, id_company, is_creater) 
		VALUES (:name, :id_company, :is_creater)
		RETURNING id_role
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRowxContext(ctx, role).Scan(&role.ID)
	if err != nil {
		return fmt.Errorf("failed to create roleincompany and get its ID: %w", err)
	}

	return nil
}

func (r *roleincompanyRepository) Update(ctx context.Context, role *models.RoleInCompany) error {
	err := role.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE Role_in_Company 
		SET name = :name, id_company = :id_company,
			is_creater = :is_creater
		WHERE id_role = :id_role
	`
	_, err = r.db.NamedExecContext(ctx, query, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (r *roleincompanyRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Role_in_Company WHERE id_role = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
