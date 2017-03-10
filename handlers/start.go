package handlers

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/RoflCopter24/citation-db/models"
	"github.com/gorilla/context"
)

func HandleStart (writer http.ResponseWriter, request *http.Request) {

	u := context.Get(request,"User")

	if u == nil {
		fmt.Print("User is nil!")
		return
	}

	user := u.(*models.User)

	pData := models.Page{ Title: "Startseite", User: user }

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(writer, "index.html", pData)
}
