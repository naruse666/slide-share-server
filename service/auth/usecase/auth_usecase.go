package usecase

import (
	"slide-share/infrastructure/repository/firebase"
	"slide-share/model"
	"time"
)

type IAuthUsecase interface {
	SignIn(user model.SignInUser) (*model.User, error)
}

type authUsecase struct {
	ur firebase.IUserRepository
}

func NewAuthUsecase(ur firebase.IUserRepository) IAuthUsecase {
	return &authUsecase{ur: ur}
}

func (au *authUsecase) SignIn(signInUser model.SignInUser) (*model.User, error) {
	storeUser, err := au.ur.GetUserByEmail(signInUser.Email)
	if err != nil {
		return nil, err
	}

	if storeUser == nil {
		user, err := au.ur.CreateUser(model.User{
			ID:        signInUser.ID,
			Name:      signInUser.Name,
			Email:     signInUser.Email,
			Image:     signInUser.Image,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return storeUser, nil
}
