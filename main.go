package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
	mgosession "github.com/RoflCopter24/citation-db/nigronimgosession"
	"github.com/RoflCopter24/negroni-sessions"
	"github.com/RoflCopter24/citation-db/settings"
	"github.com/RoflCopter24/citation-db/middleware"
	"github.com/RoflCopter24/citation-db/handlers"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"strconv"
	"github.com/RoflCopter24/negroni-sessions/mongostore"
	"github.com/gorilla/securecookie"
)

var (
	appSettings settings.AppSettings
)

func setupMgo(n *negroni.Negroni, s *settings.AppSettings) *mgo.Session {

	fmt.Println("Connecting to MongoDB: ", s.DbServer)
	fmt.Println("Database Name: ", s.DbName)

	connStr := s.GenMgoConnStr()

	// Creating the database accessor here.
	// Pointer to this database accessor will be passed to the middleware.
	dbAccessor, err := mgosession.NewDatabaseAccessor(connStr, s.DbName, "users")
	if err != nil {
		panic(err)
	}
	// Registering the middleware here.
	n.Use(mgosession.NewDatabase(*dbAccessor).Middleware())

	return dbAccessor.Clone()
}

func main() {
	appSettings = settings.AppSettings{}
	appSettings.LffOrDefault()

	n := negroni.Classic()

	// Setup MongoDb connection stuff
	s := setupMgo(n, &appSettings)


	n.Use(negroni.NewRecovery())

	n.Use(negroni.NewLogger())

	store := mongostore.New(*s, appSettings.DbName, "sessions", 900000, true, securecookie.GenerateRandomKey(16), securecookie.GenerateRandomKey(16))
	//store := cookiestore.New([]byte("citation-db.C_Store01"))
	n.Use(sessions.Sessions("CitationSession", store))

	whiteList := make([]string,5)
	whiteList[0] = "/login"
	whiteList[1] = "/css/"
	whiteList[2] = "/js/"
	whiteList[3] = "/img/"
	whiteList[4] = "/favicon.ico"

	ac := middleware.AuthChecker{ UsersDB: appSettings.DbName, UsersCollection: "users", WhiteList: whiteList }
	n.Use(ac.Middleware())

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleIndex)

	mux.HandleFunc("/start", handlers.HandleStart)

	mux.HandleFunc("/login", handlers.HandleLogin)

	mux.HandleFunc("/logout", handlers.HandleLogout)

	mux.HandleFunc("/books/add", handlers.HandleBooksAdd)

	mux.HandleFunc("/books", handlers.HandleBooksIndex)

	mux.HandleFunc("/books/", handlers.HandleBooksIndex)

	mux.HandleFunc("/books/list", handlers.HandleBooksList)
	mux.HandleFunc("/books/search", handlers.HandleBooksSearchJSON)

	mux.HandleFunc("/books/edit", handlers.HandleBooksEdit)
	mux.HandleFunc("/books/edit/", handlers.HandleBooksEdit)

	mux.HandleFunc("/books/remove/", handlers.HandleBooksDelete)

	mux.HandleFunc("/quotes/", handlers.HandleQuotesIndex)

	mux.HandleFunc("/quotes/search", handlers.HandleQuotesSearch)
	mux.HandleFunc("/quotes/search/", handlers.HandleQuotesSearch)

	mux.HandleFunc("/quotes/add", handlers.HandleQuotesAdd)

	mux.HandleFunc("/quotes/edit", handlers.HandleQuotesEdit)
	mux.HandleFunc("/quotes/edit/", handlers.HandleQuotesEdit)

	mux.HandleFunc("/quotes/remove/", handlers.HandleQuotesRemove)

	mux.Handle("/img", http.FileServer(http.Dir("img")))
	mux.Handle("/js", http.FileServer(http.Dir("js")))
	mux.Handle("/css", http.FileServer(http.Dir("css")))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	log.Println(os.Getwd())
	if appSettings.WorkingDir != "" {
		os.Chdir(appSettings.WorkingDir)
	}
	n.UseHandler(mux)
	n.Run(":" + strconv.Itoa(appSettings.Port))
}
