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
	"gopkg.in/mgo.v2"
	"log"
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
	pData.User = user

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-search.html", pData)
}

func handleQuotesSearchPOST(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "User").(*models.User)
	pData := models.PageQuoteSearchResponse{}
	pData.Title = "Suchergebnisse"
	pData.User = user

	qBook := r.PostFormValue("searchBook")
	qText := r.PostFormValue("searchText")
	qTags := r.PostFormValue("searchTags")
	qCats := r.PostFormValue("searchCategories")

	db := context.Get(r, "db").(*mgo.Database)

	if qBook != "" {
		book := models.Book{}
		errBook := db.C("books").Find(bson.M{"_id": qBook}).One(&book)

		if errBook != nil {
			pData.Error = errBook.Error()
			log.Fatal(errBook)
			goto renderHTML
		}

		for _, quote := range book.Quotations {
			if checkText(&quote, qText) || checkCats(&quote, qCats) || checkTags(&quote, qTags) {
				pData.Quotes = append(pData.Quotes, quote)
			}
		}

	} else {
		query := bson.M{ "$and": []bson.M{ {"$or" : []bson.M{
			{"$and": []bson.M{ // Username matches and group flag not set
					{"quotes.permissions.userId": user.Username },
					{"quotes.permissions.group": false },
			}},
			{"$and":[]bson.M{ // Group matches and group flag is set
						{"quotes.permissions.groupId": 0 },
						{"quotes.permissions.group": true },
			}},
		},},},
			"$or": []bson.M{
				{ "$or": []bson.M {
					{ "quotes.coreStatement": 	qText},
					{ "quotes.text": 		qText},
				}},
				{ "quotes.tags":		qTags},
				{ "quotes.categories":		qCats},
			},
		}

		//pipe := db.C("books").Pipe([]bson.M{{"$match": bson.M{"quotes.coreStatement":"John"},}})
		//resp := []bson.M{}
		//err := pipe.All(&resp)
		books := make([]models.Book, 0)
		err := db.C("books").Find(query).All(&books)

		if err != nil {
			pData.Error = err.Error()
			log.Fatal(err)
			goto renderHTML
		}

		for _, book := range books {
			pData.Quotes = append(pData.Quotes, book.Quotations...)
		}

	}

renderHTML:
	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-results.html", pData)
}

func checkText(quote *models.Quotation, text string) bool {
	if strings.Contains(quote.CoreStatement, text) {
		return true
	}
	if strings.Contains(quote.Text, text) {
		return true
	}
	return false
}

func checkTags(quote *models.Quotation, tags string) bool {
	tagsArr := strings.Split(tags, ",")
	removeWhiteSpacesInSlice(tagsArr)

	for _, tag := range tagsArr {
		for _, a := range quote.Tags {
			if strings.Contains(a, tag) {
				return true
			}
		}
	}
	return false
}

func checkCats(quote *models.Quotation, cats string) bool {
	catsArr := strings.Split(cats, ",")
	removeWhiteSpacesInSlice(catsArr)

	for _, cat := range catsArr {
		for _, a := range quote.Tags {
			if strings.Contains(a, cat) {
				return true
			}
		}
	}
	return false
}

func removeWhiteSpacesInSlice(slice []string) {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
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
	pData.User = user

	bookId := r.PostFormValue("selectedBook")

	if bookId == "" {
		pData.Error = "Kein Werk ausgewählt!"
		tpl, _ := template.ParseGlob("html/*.html")
		tpl.ExecuteTemplate(w, "quotes-add.html", pData)
		return
	}

	quote := quoteObjFromReq(r, user, "")

	db := context.Get(r, "db").(*mgo.Database)
	book := models.Book{}
	err := db.C("books").Find(bson.M{ "_id": bookId}).One(&book)

	if err != nil {
		pData.Error = "Kein werk mit dieser ID!"
	} else {
		book.Quotations = append(book.Quotations, quote)

		errI := db.C("books").UpdateId(bookId, &book)

		if errI != nil {
			panic(errI)
		}
	}

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
	pData := models.PageQuoteEdit{}
	pData.User = user

	if len(r.URL.Path) < len("/quotes/edit/a") {
		pData.Error = "Kein Werk ausgewählt!"
		pData.Title = "Fehler"
	} else {
		pData.Title = "Zitat editieren"

		quoteId := r.URL.Path[len("/quotes/edit/"):]

		db := context.Get(r, "db").(*mgo.Database)

		book := models.Book{}
		errDb := db.C("books").Find(bson.M{ "quotes._id": quoteId }).One(&book)

		if errDb != nil {
			panic(errDb)
		}

		found := false
		for _, quote := range book.Quotations {
			if quote.Id == quoteId {
				pData.Quote = quote
				found = true
			}
		}

		if !found {
			pData.Error = "Kein Werk mit dieser Id!"
		}

		pData.Book = book
		pData.GenCheckStr()

	}

	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}

func handleQuotesEditPOST(w http.ResponseWriter, r *http.Request) {
	pData := models.PageQuoteEdit{}
	pData.Title = "Zitat einpflegen"

	user := context.Get(r, "User").(*models.User)
	pData.User = user

	if len(r.URL.Path) < len("/quotes/edit/a") {
		pData.Error = "Kein Werk ausgewählt!"
		pData.Title = "Fehler"
	} else {
		pData.Title = "Zitat editieren"

		quoteId := r.URL.Path[len("/quotes/edit/"):]

		bookId := r.PostFormValue("selectedBook")

		if bookId == "" {
			pData.Error = "Kein Werk ausgewählt!"
			tpl, _ := template.ParseGlob("html/*.html")
			tpl.ExecuteTemplate(w, "quotes-add.html", pData)
			return
		}

		quote := quoteObjFromReq(r, user, quoteId)

		db := context.Get(r, "db").(*mgo.Database)
		book := models.Book{}
		err := db.C("books").Find(bson.M{"_id": bookId}).One(&book)

		if err != nil {
			pData.Error = "Kein werk mit dieser ID!"
		} else {
			//book.Quotations = append(book.Quotations, quote)

			done := false
			for i, q := range book.Quotations {
				if q.Id == quote.Id {
					book.Quotations[i] = quote
					done = true;
				}
			}

			if !done {
				pData.Error = "Zitat nicht gefunden!"
			}

			errI := db.C("books").UpdateId(bookId, &book)

			if errI != nil {
				panic(errI)
			}


		}
	}
	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}

func HandleQuotesRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user := context.Get(r, "User").(*models.User)
	pData := models.Page{}
	pData.Title = "Zitat entfernen"
	pData.User = user

	if len(r.URL.Path) < len("/quotes/edit/a") {
		pData.Error = "Kein Werk ausgewählt!"
		pData.Title = "Fehler"
	} else {
		pData.Title = "Zitat editieren"

		quoteId := r.URL.Path[len("/quotes/edit/"):]

		bookId := r.PostFormValue("selectedBook")

		if bookId == "" {
			pData.Error = "Kein Werk ausgewählt!"
			tpl, _ := template.ParseGlob("html/*.html")
			tpl.ExecuteTemplate(w, "quotes-add.html", pData)
			return
		}

		db := context.Get(r, "db").(*mgo.Database)
		book := models.Book{}
		err := db.C("books").Find(bson.M{"_id": bookId}).One(&book)

		if err != nil {
			pData.Error = "Kein werk mit dieser ID!"
		} else {
			//book.Quotations = append(book.Quotations, quote)

			done := false
			for i, q := range book.Quotations {
				if q.Id == quoteId {

					//Where a is the slice, and i is the index of the element you want to delete:
					//a = append(a[:i], a[i+1:]...)

					book.Quotations = append(book.Quotations[:i], book.Quotations[i+1:]...)
					done = true
				}
			}

			if !done {
				pData.Error = "Zitat nicht gefunden!"
			}

			errI := db.C("books").UpdateId(bookId, &book)

			if errI != nil {
				panic(errI)
			} else {
				pData.Error = "Zitat erfolgreich gelöscht"
			}


		}
	}


	tpl, _ := template.ParseGlob("html/*.html")
	tpl.ExecuteTemplate(w, "quotes-edit.html", pData)
}

func quoteObjFromReq(r *http.Request, user *models.User, id string) models.Quotation {
	quote := models.Quotation{}


	if id == "" {
		id = bson.NewObjectIdWithTime(time.Now()).Hex()
	}

	quote.Id = id
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

	return quote
}
