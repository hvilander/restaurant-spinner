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
	key := []byte(os.Getenv("SESSION_SECRET"))
	cookieStore := sessions.NewCookieStore(key)
	cookieStore.Options.HttpOnly = true
	gothic.Store = cookieStore

	// todo reconsider these options
	//store.Options.Path = "/"
	//cookieStore.MaxAge(MaxAge)
	// store.Options.Secure = IsProd // I think this breaks my current setup

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/google/callback"),
	)

}
