package usecase

import (
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
)

type IUserUsecase interface {
	UpdateUser(user model.User) (*model.User, error)
}

type userUsecase struct {
	ur firebase.IUserRepository
}

func NewUserUsecase(ur firebase.IUserRepository) IUserUsecase {
	return &userUsecase{ur: ur}
}

func (uu *userUsecase) UpdateUser(user model.User) (*model.User, error) {
	updatedUser, err := uu.ur.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
