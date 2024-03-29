package firebase

import (
	"context"
	"fmt"
	"slide-share/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
}

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) IUserRepository {
	return &UserRepository{client: client}
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	ctx := context.Background()
	iter := ur.client.Collection("users").Where("Email", "==", email).Documents(ctx)
	docSnapshot, err := iter.Next()

	if err != nil {
		if err == iterator.Done {
			// 条件に合致するドキュメントがない場合、nilを返す
			return nil, nil
		}
		// その他のエラーの場合、ログを出力してエラーを返す
		fmt.Printf("error getting user by email: %v", err)
		return nil, err
	}

	// ドキュメントが存在するかチェック
	if !docSnapshot.Exists() {
		return nil, nil
	}

	var user model.User
	if err := docSnapshot.DataTo(&user); err != nil {
		fmt.Printf("error converting user data: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) CreateUser(user model.User) (*model.User, error) {
	_, err := ur.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		fmt.Printf("error creating user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(user model.User) (*model.User, error) {
	_, err := ur.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		fmt.Printf("error updating user: %v", err)
		return nil, err
	}

	return &user, nil
}
