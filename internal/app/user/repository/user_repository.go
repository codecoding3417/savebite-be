package repository

import (
	"gorm.io/gorm"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/entity"
	"savebite/pkg/log"
)

type UserRepoItf interface {
	Create(user *entity.User) error
	Show(param dto.UserParam) (*entity.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepoItf {
	return &UserRepo{db}
}

func (r *UserRepo) Create(user *entity.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UserRepository][Create] Database error")
	}

	return err
}

func (r *UserRepo) Show(param dto.UserParam) (*entity.User, error) {
	user := entity.User{}
	err := r.db.Find(&user, param).Error
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[UserRepository][Show] Database error")
		return nil, err
	}

	return &user, nil
}
