package handler

import (
	"fmt"
	"net/http"

	layout "github.com/hvilander/restaurant-spinner/templates/layout"
)

func App(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("root accessed")

	clientID := "TEST"
	return layout.App(true, clientID).Render(r.Context(), w)

}
