package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/canonflow/canonflow-go-ddd/internal/contract"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/repository"
	"github.com/canonflow/canonflow-go-ddd/pkg/jwt"
	"github.com/canonflow/canonflow-go-ddd/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UserUsecaseImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Config         *viper.Viper
	UserRepository repository.UserRepository
	UserProducer   contract.ProducerContract
}

func NewUserUsecase(db *gorm.DB, log *logrus.Logger, config *viper.Viper, userRepository repository.UserRepository, userProducer contract.ProducerContract) UserUsecase {
	return &UserUsecaseImpl{
		DB:             db,
		Log:            log,
		Config:         config,
		UserRepository: userRepository,
		UserProducer:   userProducer,
	}
}

func (u *UserUsecaseImpl) CreateAccessToken(user *model.User) (string, error) {
	//* Generate JWT Token
	token, err := jwt.CreateToken(user, u.Config.GetString("JWT_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecaseImpl) Create(context context.Context, username string, password string) (model.User, error) {
	//* Generate Hashed Password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return model.User{}, err
	}

	//* Start DB Transaction
	tx := u.DB.WithContext(context).Begin()
	if tx.Error != nil {
		return model.User{}, tx.Error
	}
	defer tx.Rollback()

	//* Create User
	user := model.User{
		Username: username,
		Password: hashedPassword,
	}

	err = u.UserRepository.Create(tx, &user)
	if err != nil {
		return model.User{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.User{}, err
	}

	if u.UserProducer != nil {
		event := model.UserEvent{
			ID:        strconv.Itoa(int(user.ID)),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.UnixMicro(),
			UpdatedAt: user.UpdatedAt.UnixMicro(),
		}
		u.Log.Info("Publishing user created event")
		if err = u.UserProducer.Send(&event); err != nil {
			u.Log.Warnf("Failed publish user created event : %+v", err)
			return user, err
		}
	} else {
		u.Log.Info("Kafka producer is disabled, skipping user created event")
	}

	return user, nil
}

func (u *UserUsecaseImpl) Login(user *model.User, password string) error {
	//* Check Password
	if !utils.CheckPassword(password, user.Password) {
		return errors.New("Wrong Credentials")
	}

	if u.UserProducer != nil {
		event := model.UserEvent{
			ID:        strconv.Itoa(int(user.ID)),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.UnixMicro(),
			UpdatedAt: user.UpdatedAt.UnixMicro(),
		}
		u.Log.Info("Publishing user created event")
		if err := u.UserProducer.Send(&event); err != nil {
			u.Log.Warnf("Failed publish user created event : %+v", err)
			return nil
		}
	} else {
		u.Log.Info("Kafka producer is disabled, skipping user created event")
	}

	return nil
}

func (u *UserUsecaseImpl) FindByUsername(username string) (*model.User, error) {
	user, err := u.UserRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecaseImpl) FindById(id int64) (*model.User, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
