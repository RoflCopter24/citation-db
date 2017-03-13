package handlers

import (
	"net/http"
	"html/template"
	"github.com/RoflCopter24/citation-db/models"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strings"
	"strconv"
)

func HandleQuotesIndex(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "User")
	if u == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	http.Redirect(w, r, "/quotes/search", 302)
}

func HandleQuotesSearch(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		handleQuotesSearchGET(w, r)
	} else if r.Method == "POST" {
		handleQuotesSearchPOST(w,r)
	}
}

func handleQuotesSearchGET(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Zitate durchsuchen"
	pData.User = user;

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-search.html", pData)
}

func handleQuotesSearchPOST(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Suchergebnisse"
	pData.User = user;

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-results.html", pData)
}

func HandleQuotesAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleQuotesAddGET(w, r)
	} else if r.Method == "POST" {
		handleQuotesAddPOST(w,r)
	}
}

func handleQuotesAddGET(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)

	pData := models.PageQuoteEdit{}
	pData.Title = "Zitat einpflegen"
	pData.User = user;

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-add.html", pData)
}

func handleQuotesAddPOST(w http.ResponseWriter, r *http.Request) {
	pData := models.PageQuoteEdit{}
	pData.Title = "Zitat einpflegen"

	user := context.Get(r, "User").(*models.User)
	pData.User = user;

	quote := models.Quotation{}

	quote.Id = bson.NewObjectIdWithTime(time.Now()).Hex()
	quote.Categories = strings.Split(r.PostFormValue("quoteCategories"), ",")
	quote.CoreStatement = r.PostFormValue("quoteCoreStatement")
	quote.CreationDate = time.Now()
	quote.CreatorId = user.Username
	quote.Description = r.PostFormValue("quoteDescription")
	quote.Files = make([]string, 0)
	quote.History = make([]models.Quotation, 0)
	quote.LastModified = time.Now()
	quote.Legend = r.PostFormValue("quoteLegend")

	pNrS, errPNrS := strconv.Atoi(r.PostFormValue("quotePageNrStart"))
	if errPNrS != nil {
		pNrS = 1
	}
	quote.PageNrStart = pNrS

	pNrSt, errPNrSt := strconv.Atoi(r.PostFormValue("quotePageNrStop"))
	if errPNrSt != nil {
		pNrSt = 1
	}
	quote.PageNrStop = pNrSt
	quote.Permissions = make([]models.Permission, 1)
	quote.Permissions[0] = models.Permission{UserId:user.Username, Read: true, Write: true}
	quote.Tags = strings.Split(r.PostFormValue("quoteTags"), ",")
	quote.Text = r.PostFormValue("quoteText")

	t, errT := strconv.Atoi(r.PostFormValue("quoteType"))
	if errT != nil {
		t = 0
	}
	quote.Type = t

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-add.html", pData)
}

func HandleQuotesEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleQuotesEditGET(w, r)
	} else if r.Method == "POST" {
		handleQuotesEditPOST(w,r)
	}
}

func handleQuotesEditGET(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Zitat bearbeiten"
	pData.User = user

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}

func handleQuotesEditPOST(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Zitat bearbeiten"
	pData.User = user

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}

func HandleQuotesRemove(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Zitat entfernen"
	pData.User = user

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}
