package usecase

import (
	"context"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
)

type UserUsecase interface {
	CreateAccessToken(user *model.User) (string, error)
	Create(ctx context.Context, username string, password string) (model.User, error)
	Login(user *model.User, password string) error
	FindByUsername(username string) (*model.User, error)
}
