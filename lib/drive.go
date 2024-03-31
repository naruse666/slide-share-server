package lib

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func InitDrive() (*drive.Service, error) {
	ctx := context.Background()

	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	credentials, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("FIREBASE_ACCOUNT")), "https://www.googleapis.com/auth/drive")
	if err != nil {
		log.Printf("error credentials from json: %v\n", err)
	}

	opt := option.WithCredentials(credentials)

	driveService, err := drive.NewService(ctx, opt)
	if err != nil {
		log.Fatalf("error initializing drive client: %v", err)
	}

	return driveService, nil
}
