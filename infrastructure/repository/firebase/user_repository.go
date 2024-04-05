package firebase

import (
	"context"
	"fmt"
	"slide-share/model"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type IUserRepository interface {
	GetUser(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUsers() ([]model.User, error)
	CreateUser(user model.User) (*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
}

type UserRepository struct {
	client *firestore.Client
}

func NewUserRepository(client *firestore.Client) IUserRepository {
	return &UserRepository{client: client}
}

func (ur *UserRepository) GetUser(id string) (*model.User, error) {
	ctx := context.Background()
	docSnapshot, err := ur.client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		fmt.Printf("error getting user by id: %v", err)
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

func (ur *UserRepository) GetUsers() ([]model.User, error) {
	ctx := context.Background()
	iter := ur.client.Collection("users").OrderBy("Role", firestore.Asc).OrderBy("CreatedAt", firestore.Asc).Documents(ctx)

	var users []model.User
	for {
		docSnapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("error getting users: %v", err)
			return nil, err
		}

		var user model.User
		if err := docSnapshot.DataTo(&user); err != nil {
			fmt.Printf("error converting user data: %v", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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
	user.UpdatedAt = time.Now()
	_, err := ur.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	if err != nil {
		fmt.Printf("error updating user: %v", err)
		return nil, err
	}

	return &user, nil
}
