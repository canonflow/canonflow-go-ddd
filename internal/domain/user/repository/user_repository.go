package repository

import (
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(db *gorm.DB, user *model.User) error
	Update(db *gorm.DB, user *model.User) error
	FindByID(id int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}
