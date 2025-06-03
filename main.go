package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/hvilander/restaurant-spinner/handler"
	"github.com/hvilander/restaurant-spinner/internal/auth"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func main() {

	// setup config
	godotenv.Load()
	PORT := os.Getenv("PORT")
	sessSec := os.Getenv("SESSION_SECRET")
	fmt.Println("Starting Up...")
	fmt.Println("sess sec", sessSec)

	auth.NewAuth()

	m := map[string]string{
		"google": "Google",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		fmt.Println("in callback")

		fmt.Println("query", q[":provider"])
		fmt.Println(gothic.SessionName)
		sesh, seshErr := gothic.Store.Get(req, gothic.SessionName)
		if seshErr != nil {
			fmt.Println("sesh error:", seshErr)
			fmt.Println("sesh error:", seshErr)
		}

		fmt.Println("len", len(sesh.Values))
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)

	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		// try to get the user without re-authenticating
		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
			fmt.Println("calling complete user auth")
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(res, gothUser)
		} else {
			fmt.Println("calling begin")
			gothic.BeginAuthHandler(res, req)
		}
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, providerIndex)
	})

	log.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", p))

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	// maybe set filepathroot to an env var
	filepathRoot := "./app"
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", appHandler)

	// register paths
	mux.Handle("/app/home", handler.MakeHandler(handler.HandlerHomeIndex))

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
