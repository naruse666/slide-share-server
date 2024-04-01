package lib

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"github.com/joho/godotenv"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func InitFirebase() (*firestore.Client, *storage.Client) {
	ctx := context.Background()

	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	credentials, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("FIREBASE_ACCOUNT")), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("error credentials from json: %v\n", err)
	}

	opt := option.WithCredentials(credentials)
	config := &firebase.Config{
		StorageBucket: os.Getenv("STORAGE_BUCKET"),
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing firestore client: %v", err)
	}

	storage, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalf("error initializing storage client: %v", err)
	}

	return firestore, storage
}
