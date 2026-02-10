package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type CompanyRepository interface {
	GetAll(ctx context.Context) ([]models.Company, error)
	GetByID(ctx context.Context, id int) (*models.Company, error)
	Create(ctx context.Context, company *models.Company) error
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id int) error
}

type companyRepository struct {
	db QueryExecutor
}

func NewCompanyRepository(db QueryExecutor) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetAll(ctx context.Context) ([]models.Company, error) {
	var companies []models.Company
	query := `SELECT * FROM companies`
	err := r.db.SelectContext(ctx, &companies, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get companies: %w", err)
	}
	return companies, nil
}
func (r *companyRepository) GetByID(ctx context.Context, id int) (*models.Company, error) {
	var company models.Company

	query := `SELECT * FROM companies WHERE id_company = $1`

	err := r.db.GetContext(ctx, &company, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}
	return &company, nil
}
func (r *companyRepository) Create(ctx context.Context, company *models.Company) error {
	err := company.Validate()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO companies (name, description, email, address, inn, egrul) 
		VALUES (:name, :description, :email, :address, :inn, :egrul)
		RETURNING id_company
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to prepare query: %w", err)
    }
    defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, company).Scan(&company.ID)
	if err != nil {
		return fmt.Errorf("failed to create company and get its ID: %w", err)
	}

	return nil
}
func (r *companyRepository) Update(ctx context.Context, company *models.Company) error {
	err := company.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE companies 
		SET name = :name, description = :description, 
			email = :email, address = :address, 
			inn = :inn, egrul = :egrul
		WHERE id_company = :id_company
	`
	_, err = r.db.NamedExecContext(ctx, query, company)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}
func (r *companyRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM companies WHERE id_company = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}
	return nil
}
