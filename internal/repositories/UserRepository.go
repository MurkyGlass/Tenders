package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db QueryExecutor
}

func NewUserRepository(db QueryExecutor) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users`
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE id_user = $1`
	
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE login = $1`
	
	err := r.db.GetContext(ctx, &user, query, login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	query := `
		INSERT INTO users (login, name, email, password, id_company, id_role_in_company, id_role) 
		VALUES (:login, :name, :email, :password, :id_company, :id_role_in_company, :id_role)
		RETURNING id_user
	`
	namedQuery, err := r.db.PrepareNamedContext(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to prepare query(user): %w", err)
    }
    defer namedQuery.Close()
    
    err = namedQuery.QueryRowContext(ctx, user).Scan(&user.ID)
    if err != nil {
        return fmt.Errorf("failed to create user and get its ID: %w", err)
    }

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}
	//password not updates
	query := `
		UPDATE users 
		SET login = :login, name = :name, email = :email, 
			id_company = :id_company,
			id_role_in_company = :id_role_in_company,
			id_role = :id_role
		WHERE id_user = :id_user
	`
	_, err = r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id_user = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
