package repositories

import (
	"context"
	"fmt"
	"main/internal/repositories/models"
)

type LogRepository interface {
	GetAll(ctx context.Context) ([]models.Log, error)
	GetByID(ctx context.Context, id int) (*models.Log, error)
	Create(ctx context.Context, log *models.Log) LinkerLog
	Update(ctx context.Context, log *models.Log) error
	Delete(ctx context.Context, id int) error
}

type logRepository struct {
	db QueryExecutor
}

func NewLogRepository(db QueryExecutor) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) GetAll(ctx context.Context) ([]models.Log, error) {
	var logs []models.Log
	query := `SELECT * FROM logs`
	err := r.db.SelectContext(ctx, &logs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, nil
}

func (r *logRepository) GetByID(ctx context.Context, id int) (*models.Log, error) {
	var log models.Log

	query := `SELECT * FROM logs WHERE id_log = $1`

	err := r.db.GetContext(ctx, &log, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}
	return &log, nil
}
func (r *logRepository) Create(ctx context.Context, log *models.Log) LinkerLog {
	err := log.Validate()
	if err != nil {
		return &linkerLog{err: err}
	}
	query := `
		INSERT INTO logs (id_user, id_entity, id_type, datetime_create) 
		VALUES (:id_user, :id_entity, :id_type, :datetime_create)
		RETURNING id_log
	`
	err = r.db.QueryRowxContext(ctx, query, log).Scan(&log.ID)
	if err != nil {
		return &linkerLog{err: fmt.Errorf("failed to create log: %w", err)}
	}

	return &linkerLog{idLog: log.ID, db: r.db, err: nil}
}

func (r *logRepository) Update(ctx context.Context, log *models.Log) error {
	err := log.Validate()
	if err != nil {
		return err
	}
	query := `
		UPDATE logs
		SET id_user = :id_user, id_entity = :id_entity,
			id_type = :id_type, datetime_create = :datetime_create
		WHERE id_log = :id_log
	`
	_, err = r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}
	return nil
}

func (r *logRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM logs WHERE id_log = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete log: %w", err)
	}
	return nil
}
