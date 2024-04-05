package usecase

import (
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
)

type IUserUsecase interface {
	GetUser(id string) (*model.User, error)
	GetUsers() ([]model.User, error)
	UpdateUser(user model.User) (*model.User, error)
}

type userUsecase struct {
	ur firebase.IUserRepository
}

func NewUserUsecase(ur firebase.IUserRepository) IUserUsecase {
	return &userUsecase{ur: ur}
}

func (uu *userUsecase) GetUser(id string) (*model.User, error) {
	user, err := uu.ur.GetUser(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *userUsecase) GetUsers() ([]model.User, error) {
	users, err := uu.ur.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uu *userUsecase) UpdateUser(user model.User) (*model.User, error) {
	updatedUser, err := uu.ur.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
