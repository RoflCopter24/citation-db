package handlers

import (
	"net/http"
	"github.com/gorilla/context"
)

func HandleIndex (w http.ResponseWriter, r *http.Request) {

	u := context.Get(r, "User")
	if u == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	http.Redirect(w, r, "/start", 302)
}
