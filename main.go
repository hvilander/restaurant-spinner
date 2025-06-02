package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/hvilander/restaurant-spinner/handler"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Up...")

	// setup config
	godotenv.Load()
	PORT := os.Getenv("PORT")

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
