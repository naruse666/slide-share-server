package firebase

import (
	"context"
	"log"
	"slide-share/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type IUserRepository interface {
	GetUsers() ([]model.User, error)
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

func (ur *UserRepository) GetUsers() ([]model.User, error) {
	users, err := ur.client.Collection("users").Documents(context.Background()).GetAll()
	if err != nil {
		log.Fatalf("error getting user collection: %v", err)
	}

	var userCollection []model.User
	for _, user := range users {
		var u model.User
		user.DataTo(&u)
		userCollection = append(userCollection, u)
	}

	return userCollection, nil
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
		log.Printf("error getting user by email: %v", err)
		return nil, err
	}

	// ドキュメントが存在するかチェック
	if !docSnapshot.Exists() {
		return nil, nil
	}

	var user model.User
	if err := docSnapshot.DataTo(&user); err != nil {
		log.Printf("error converting user data: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) CreateUser(user model.User) (*model.User, error) {
	_, err := ur.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		log.Fatalf("error creating user: %v", err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(user model.User) (*model.User, error) {
	_, err := ur.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		log.Fatalf("error updating user: %v", err)
		return nil, err
	}

	return &user, nil
}
