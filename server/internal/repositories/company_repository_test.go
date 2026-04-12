package repositories_test

import (
	"context"
	"main/internal/repositories"
	"main/internal/repositories/models"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
    db, err := sqlx.Connect("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to connect to in-memory DB: %v", err)
    }

    schema := `
    CREATE TABLE companies (
        id_company INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(50) NOT NULL UNIQUE,
        address VARCHAR(100) NOT NULL,
        inn VARCHAR(12) NOT NULL UNIQUE,
        egrul VARCHAR(13) NOT NULL UNIQUE,
        description VARCHAR(500)
    );`
    _, err = db.Exec(schema)
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }
    return db
}

func TestCompanyRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := repositories.NewCompanyRepository(db)

    company := &models.Company{
        Name:        "ООО Тест",
        Email:       "test@test.ru",
        Address:     "Москва, ул. Тестовая, 1",
        INN:         "1234567890",
        EGRUL:       "1234567890123",
        Description: "Тестовая компания",
    }

    err := repo.Create(context.Background(), company)
    if err != nil {
        t.Errorf("Create() error = %v, want nil", err)
    }
    if company.ID == 0 {
        t.Error("Company.ID not set after Create")
    }

    var count int
    err = db.Get(&count, "SELECT COUNT(*) FROM companies WHERE id_company = ?", company.ID)
    if err != nil {
        t.Errorf("Failed to query DB: %v", err)
    }
    if count != 1 {
        t.Errorf("Company not found in DB, count = %d", count)
    }
}