package auth

import (
	"context"

	"chat-app/config"
	authRepo "chat-app/internal/api/repository/auth"
	"chat-app/internal/dto"
	"chat-app/internal/models"
	"chat-app/pkg/customError"
	"chat-app/pkg/logger"
	"chat-app/pkg/tracer"
	"chat-app/pkg/utils"

	"github.com/pkg/errors"
)

type Service interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, eamil, password string) (*dto.LoginResDto, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type authSrv struct {
	cfg      *config.Config
	authRepo authRepo.Repository
	logger   logger.Logger
}

func NewService(cfg *config.Config, authRepo authRepo.Repository, logger logger.Logger) Service {

	return &authSrv{
		cfg:      cfg,
		authRepo: authRepo,
		logger:   logger,
	}
}

func (srv *authSrv) Register(ctx context.Context, user *models.User) error {
	c, span := tracer.NewSpan(ctx, "AuthService.Register", nil)
	defer span.End()

	// hashing password before insert into database
	hashPassword, err := utils.HashText(user.Password)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "AuthService.Register")
	}
	user.Password = hashPassword

	if err := srv.authRepo.Create(c, user); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "AuthService.Register")
	}

	return nil
}

func (srv *authSrv) Login(ctx context.Context, email string, password string) (*dto.LoginResDto, error) {
	ctx, span := tracer.NewSpan(ctx, "AuthService.Login", nil)
	defer span.End()

	user, err := srv.authRepo.GetByEmail(ctx, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "AuthService.Login")
	}

	if err := utils.CheckTextHash(password, user.Password); err != nil {
		tracer.AddSpanError(span, err)
		return nil, customError.ErrPasswordCodeNotMatched
	}
	accessJWT, err := utils.GetJwt(srv.cfg, user.ID, user.Username, user.Email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, customError.ErrInternalServerError
	}

	return &dto.LoginResDto{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		Token:     accessJWT,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (srv *authSrv) GetAll(ctx context.Context) ([]*models.User, error) {

	ctx, span := tracer.NewSpan(ctx, "AuthService.GetAll", nil)
	defer span.End()

	return srv.authRepo.GetAll(ctx)
}

func (srv *authSrv) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "AuthService.GetByEmail", nil)
	defer span.End()

	user, err := srv.authRepo.GetByEmail(ctx, email)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "AuthService.GetByEmail")
	}
	return user, nil
}
