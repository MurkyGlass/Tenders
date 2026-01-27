package repositories

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
    ErrNotFound    = errors.New("record not found")
    ErrInvalidInput = errors.New("invalid input data")
    ErrConflict     = errors.New("data conflict")
    ErrNotAllowed   = errors.New("operation not allowed")
)

type Repository interface {
	Users() UserRepository


	BeginTx(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error

	Users() UserRepository
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Users() UserRepository {
	return NewUserRepository(r.db)
}

func (r *repository) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

// Реализация транзакции
type transaction struct {
	tx *sqlx.Tx
}


func (t *transaction) Commit() error {
	return t.tx.Commit()
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
func (t *transaction) Users() UserRepository {
	return NewUserRepository(t.tx)
}

