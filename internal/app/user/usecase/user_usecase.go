package usecase

import (
	"github.com/google/uuid"
	"savebite/internal/app/user/repository"
	"savebite/internal/domain/dto"
)

type UserUsecaseItf interface {
	GetProfile(id uuid.UUID) (dto.UserProfile, error)
}

type UserUsecase struct {
	userRepo repository.UserRepoItf
}

func NewUserUsecase(userRepo repository.UserRepoItf) UserUsecaseItf {
	return &UserUsecase{userRepo}
}

func (u *UserUsecase) GetProfile(id uuid.UUID) (dto.UserProfile, error) {
	param := dto.UserParam{ID: id}
	user, err := u.userRepo.Show(param)
	if err != nil {
		return dto.UserProfile{}, err
	}

	return dto.UserProfile{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
