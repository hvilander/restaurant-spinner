package auth

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	MaxAge = 86400 * 30 // 30 days
	IsProd = false
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_KEY")
	googleClientSecret := os.Getenv("GOOGLE_SECRET")
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/callback?provider=google"),
	)

}
