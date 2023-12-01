package auth

import (
	"context"

	"chat-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type authRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &authRepo{db: db}
}
