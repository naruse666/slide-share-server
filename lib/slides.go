package lib

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

func InitSlides() (*slides.Service, error) {
	ctx := context.Background()

	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}

	credentials, err := google.CredentialsFromJSON(ctx, []byte(os.Getenv("FIREBASE_ACCOUNT")), "https://www.googleapis.com/auth/presentations")
	if err != nil {
		log.Printf("error credentials from json: %v\n", err)
	}

	opt := option.WithCredentials(credentials)

	slidesService, err := slides.NewService(ctx, opt)
	if err != nil {
		log.Fatalf("error initializing slides client: %v", err)
	}

	return slidesService, nil
}
