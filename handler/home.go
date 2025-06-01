package handler

import (
	"fmt"
	"net/http"

	home "github.com/hvilander/restaurant-spinner/templates/home"
)

func HandlerHomeIndex(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("home accessed")

	return home.Index().Render(r.Context(), w)

}
