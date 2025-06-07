package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hvilander/restaurant-spinner/handler"
	"github.com/hvilander/restaurant-spinner/internal/auth"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "restaurant-spinner-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)

}

func main() {
	// load env vars
	godotenv.Load()
	PORT := os.Getenv("PORT")

	// init goth / gothic auth
	auth.NewAuth()

	m := map[string]string{
		"google": "Google",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	// maybe set filepathroot to an env var
	//filepathRoot := "./app"
	//appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	//mux.Handle("/app/", appHandler)

	// register paths
	//mux.Handle("/app/home", handler.MakeHandler(handler.HandlerHomeIndex))

	rootHandler := handler.MakeHandler(func(res http.ResponseWriter, req *http.Request) error {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, providerIndex)
		return nil
	})

	authProviderHandler := handler.MakeHandler(func(res http.ResponseWriter, req *http.Request) error {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(res, gothUser)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
		return nil
	})

	authCallbackHandler := handler.MakeHandler(func(res http.ResponseWriter, req *http.Request) error {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return err
		}
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
		// make a jwt
		dummyUserId := uuid.New()
		token, err := MakeJWT(dummyUserId, "dummy-secret", time.Minute*3)
		fmt.Println(token)

		return nil

	})

	logoutHandler := handler.MakeHandler(func(res http.ResponseWriter, req *http.Request) error {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
		return nil
	})

	mux.Handle("GET /", rootHandler)
	mux.Handle("GET /auth/{provider}", authProviderHandler)
	mux.Handle("GET /auth/{provider}/callback", authCallbackHandler)
	mux.Handle("GET /auth/{provider}/logout", logoutHandler)
	mux.Handle("GET /test", handler.MakeHandler(
		func(res http.ResponseWriter, req *http.Request) error {
			authHeader := res.Header().Get("Authorization")
			fmt.Println("authHeader:", authHeader)

			return nil
		}))

	// start server
	slog.Info(fmt.Sprintf("server starting on http://localhost:%s", PORT))
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("error starting up", "error", err)
	}

}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`
