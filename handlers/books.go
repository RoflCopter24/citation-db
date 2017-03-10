package handlers

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/RoflCopter24/citation-db/models"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
	"strconv"
	"github.com/goincremental/negroni-sessions"
	"gopkg.in/mgo.v2/bson"
)

func HandleBooksAdd (writer http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		HandleBooksAddGET(writer, request)
	} else if request.Method == "POST" {
		HandleBooksAddPOST(writer, request)
	}
}

func HandleBooksAddGET (w http.ResponseWriter, r *http.Request) {

	u := context.Get(r,"User")

	if u == nil {
		fmt.Print("User is nil!")
		return
	}

	user := u.(*models.User)

	pData := models.Page{ Title: "Werk hinzufügen", User: user }
	pData.GenCheckStr()

	session := sessions.GetSession(r)
	session.Set("CheckStr", pData.CheckStr)

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "books-add.html", pData)
}

func HandleBooksAddPOST (w http.ResponseWriter, r *http.Request) {

	session := sessions.GetSession(r)
	checkStr := session.Get("CheckStr").(string)

	if checkStr != r.PostFormValue("checkField") {
		fmt.Println("Checkfield did not match! Expected: " + checkStr + ", got: " + r.PostFormValue("checkField"))
		http.Error(w, "Access denied. Form abuse detected.", http.StatusForbidden)
	}

	session.Delete("CheckStr")
	// You can access the mgo db object from the request object.
	// The db object is stored in key `db`.
	db := context.Get(r, "db").(*mgo.Database)
	// Now lets perform a count query using mgo db object.

	fPublishedDate, errFP := time.Parse("yyyy-MM-dd", r.PostFormValue("booksFirstPublished"))

	if errFP != nil {
		fmt.Println("[AddBook] Published Date failed to convert -> " + errFP.Error())
		fPublishedDate = time.Date(1970, 1,1, 0,0,0,0, time.Local)
	}

	lCheckedDate, errLC := time.Parse("yyyy-MM-dd", r.PostFormValue("bookLastChecked"))

	if errLC != nil {
		fmt.Println("[AddBook] Last Checked date failed to convert -> " + errLC.Error())
		lCheckedDate = time.Now()
	}

	edition, errEdition := strconv.Atoi(r.PostFormValue("booksEdition"))

	if errEdition != nil {
		fmt.Println("[AddBook] Edition failed to convert -> " + errEdition.Error())
		edition = 1
	}

	isbn, errISBN := strconv.Atoi(r.PostFormValue("bookISBN"))

	if errISBN != nil {
		fmt.Println("[AddBook] ISBN failed to convert -> " + errISBN.Error())
		isbn = 0
	}

	lineBookNr, errLBN := strconv.Atoi(r.PostFormValue("bookLineNumber"))

	if errLBN != nil {
		fmt.Println("[AddBook] LineBookNr failed to convert -> " + errLBN.Error())
		lineBookNr = 0
	}

	numOfB, errNOB := strconv.Atoi(r.PostFormValue("bookNumberOfBooks"))

	if errNOB != nil {
		fmt.Println("[AddBook] numOfB failed to convert -> " + errNOB.Error())
		numOfB = 1
	}

	price, errPrice := strconv.ParseFloat(r.PostFormValue("bookPrice"), 64)

	if errPrice != nil {
		fmt.Println("[AddBook] Price failed to convert -> " + errPrice.Error())
		price = 0.0
	}

	rating, errRating := strconv.Atoi(r.PostFormValue("booksRating"))

	if errRating != nil {
		fmt.Println("[AddBook] Rating failed to convert -> " + errRating.Error())
		rating = 0
	}

	t, errT := strconv.Atoi(r.PostFormValue("bookType"))

	if errT != nil {
		fmt.Println("[AddBook] Book Type failed to convert -> " + errT.Error())
		t = 0
	}

	year, errYear := strconv.Atoi(r.PostFormValue("bookYear"))

	if errYear != nil {
		fmt.Println("[AddBook] Year failed to convert -> " + errYear.Error())
		year = 1970
	}

	b := models.Book{
		Title: 			r.PostFormValue("bookTitle"),
		Abstract: 		r.PostFormValue("bookAbstract"),
		AdditionalUrl: 		r.PostFormValue("bookAddUrl"),
		AddTitle: 		r.PostFormValue("bookAddTitle"),
		AuthoringAssistants: 	r.PostFormValue("bookAuthoringAssistants"),
		Authors: 		r.PostFormValue("bookAuthor"),
		Categories:		strings.Split(r.PostFormValue("bookCategories"), ","),
		DateAdded:		time.Now(),
		DateModified:		time.Now(),
		Edition:		edition,
		FirstPublished:		fPublishedDate,
		Institution:		r.PostFormValue("bookInstitution"),
		ISBN:			isbn,
		Language:		r.PostFormValue("bookLanguage"),
		LastChecked:		lCheckedDate,
		LineBookNo:		lineBookNr,
		LineTitle:		r.PostFormValue("bookLineTitle"),
		Notes:			r.PostFormValue("bookNotes"),
		NumberOfBooks:		numOfB,
		Price:			price,
		Publisher:		r.PostFormValue("bookPublisher"),
		PublishingLocation:	r.PostFormValue("bookPublishingLocation"),
		Locations:		make([]models.Location, 0),
		Quotations:		make([]models.Quotation, 0),
		Rating:			rating,
		References:		strings.Split(r.PostFormValue("bookReferences"), ","),
		SubTitle:		r.PostFormValue("bookSubtitle"),
		Tags:			strings.Split(r.PostFormValue("bookTags"), ","),
		Tasks:			make([]models.Task, 0),
		ToC:			r.PostFormValue("bookToc"),
		Type:			t,
		Year:			year,
	}
	b.Id = bson.NewObjectId().Hex()
	err := db.C("books").Insert(&b)

	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/books/list", http.StatusTemporaryRedirect)
}

func HandleBooksList (writer http.ResponseWriter, request *http.Request) {

	u := context.Get(request,"User")

	if u == nil {
		fmt.Print("User is nil!")
		return
	}

	// Receive current user from context
	user := u.(*models.User)

	// Determine if number of items per page is set
	numPerPage, errNPP := strconv.Atoi(request.FormValue("itemsPerPage"))
	if errNPP != nil {
		// if not limit to 30 items per page
		numPerPage = 30
	}

	db := context.Get(request, "db").(*mgo.Database)

	// See how many total books are in the Database
	numOfBooks, errCount := db.C("books").Count()
	if errCount != nil {
		panic(errCount)
	}

	// Calc how many pages are needed
	pageCount := numOfBooks / numPerPage
	var currPageStr string
	if len(request.URL.Path) > len("/books/list") {
		currPageStr = request.URL.Path[len("/books/list/"):]
	}

	if currPageStr == "" {
		currPageStr = "1"
	}

	currPage, errCurrPage := strconv.Atoi(currPageStr)

	if errCurrPage != nil {
		currPage = 1
	}

	var books []models.Book
	err := db.C("books").Find(nil).Skip((currPage-1)*numPerPage).Limit(25).All(&books)

	pData := models.PageBookList{}
	pData.Title 	= "Werke"
	pData.User 	= user
	pData.Books 	= books
	pData.PageCount = pageCount

	if err != nil {
		fmt.Println("[ListBooks] Query failed: " + err.Error())
		pData.Error = err.Error()
	}

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(writer, "books-list.html", pData)
}

func HandleBooksEditGET (writer http.ResponseWriter, request *http.Request) {

}

func HandleBooksIndex (w http.ResponseWriter, r *http.Request) {

	u := context.Get(r, "User")
	if u == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	http.Redirect(w, r, "/books/list", 302)
}
