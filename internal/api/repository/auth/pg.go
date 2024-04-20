package auth

import (
	"context"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/pkg/errors"
)

func (r *authRepo) Create(ctx context.Context, user *models.User) error {
	ctx, span := tracer.NewSpan(ctx, "AuthRepo.Create", nil)
	defer span.End()

	sqlQuery := `INSERT INTO users(email, password, username) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, sqlQuery, user.Email, user.Password, user.Username)

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "AuthRepo.Create")
	}
	return nil
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "AuthRepo.GetByEmail", nil)
	defer span.End()

	user := models.User{}
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "AuthRepo.GetByEmail")
	}

	return &user, nil
}

func (r *authRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "AuthRepo.GetAll", nil)
	defer span.End()

	users := []*models.User{}
	// if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
	// 	tracer.AddSpanError(span, err)
	// 	return nil, errors.Wrap(err, "AuthRepo.GetAll")
	// }

	return users, nil
}
