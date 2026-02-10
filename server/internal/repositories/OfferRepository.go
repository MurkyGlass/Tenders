package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type OfferRepository interface {
	GetAll(ctx context.Context) ([]models.Offer, error)
	GetByID(ctx context.Context, id int) (*models.Offer, error)
	Create(ctx context.Context, offer *models.Offer) error
	Update(ctx context.Context, offer *models.Offer) error
	Delete(ctx context.Context, id int) error
}

type offerRepository struct {
	db QueryExecutor
}

func NewOfferRepository(db QueryExecutor) OfferRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) GetAll(ctx context.Context) ([]models.Offer, error) {
	var offers []models.Offer
	query := `SELECT * FROM offers`
	err := r.db.SelectContext(ctx, &offers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get offers: %w", err)
	}
	return offers, nil
}
func (r *offerRepository) GetByID(ctx context.Context, id int) (*models.Offer, error) {
	var offer models.Offer

	query := `SELECT * FROM offers WHERE id_offer = $1`

	err := r.db.GetContext(ctx, &offer, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer: %w", err)
	}
	return &offer, nil
}
func (r *offerRepository) Create(ctx context.Context, offer *models.Offer) error {
	err := offer.Validate()
	if err != nil {
		return err
	}
	query := `
		INSERT INTO offers (name, description, datetime_create, price, id_company, id_status, id_tender) 
		VALUES (:name, :description, :datetime_create, :price, :id_company, :id_status, :id_tender)
		RETURNING id_offer
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()
	err = stmt.QueryRowxContext(ctx, offer).Scan(&offer.ID)
	if err != nil {
		return fmt.Errorf("failed to create offer and get its ID: %w", err)
	}

	return nil
}
func (r *offerRepository) Update(ctx context.Context, offer *models.Offer) error {
	err := offer.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE offers 
		SET name = :name, description = :description, 
			datetime_create = :datetime_create, price = :price, 
			id_company = :id_company, id_status = :id_status,
			id_tender = :id_tender
		WHERE id_offer = :id_offer
	`
	_, err = r.db.NamedExecContext(ctx, query, offer)
	if err != nil {
		return fmt.Errorf("failed to update offer: %w", err)
	}
	return nil
}
func (r *offerRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM offers WHERE id_offer = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete offer: %w", err)
	}
	return nil
}
