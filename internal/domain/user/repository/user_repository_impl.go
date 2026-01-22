package repository

import (
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepositoryImpl(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: DB}
}

func (repo *UserRepositoryImpl) Create(db *gorm.DB, user *model.User) error {
	return db.Create(user).Error
}

func (repo *UserRepositoryImpl) Update(db *gorm.DB, user *model.User) error {
	return db.Save(user).Error
}

func (repo *UserRepositoryImpl) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := repo.DB.First(&user, id).Error
	return &user, err
}

func (repo *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
