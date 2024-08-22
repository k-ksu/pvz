package postgres

import "HomeWork_1/internal/pkg/database"

type Repository struct {
	db *database.Database
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{db: db}
}
