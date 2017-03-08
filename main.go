package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
	mgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/RoflCopter24/citation-db/settings"
	"github.com/RoflCopter24/citation-db/middleware"
	"html/template"
	"github.com/RoflCopter24/citation-db/models"
	"github.com/RoflCopter24/citation-db/handlers"
)

var (
	appSettings settings.AppSettings
)

func setupMgo(n *negroni.Negroni, s *settings.AppSettings) {

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
}

func main() {
	appSettings = settings.AppSettings{}
	appSettings.LffOrDefault()

	n := negroni.Classic()


	n.Use(negroni.NewRecovery())

	n.Use(negroni.NewLogger())

	store := cookiestore.New([]byte("citation-db.C_Store01"))
	n.Use(sessions.Sessions("CitationSession", store))

	// Setup MongoDb connection stuff
	setupMgo(n, &appSettings)

	whiteList := make([]string,5)
	whiteList[0] = "/login"
	whiteList[1] = "/css/"
	whiteList[2] = "/js/"
	whiteList[3] = "/img/"

	ac := middleware.AuthChecker{ UsersDB: appSettings.DbName, UsersCollection: "users", WhiteList: whiteList }
	n.Use(ac.Middleware())

	mux := http.NewServeMux()

	mux.HandleFunc("/index.html", func(writer http.ResponseWriter, request *http.Request) {

		fmt.Print("Das Lama rennt!")

		sess := sessions.GetSession(request)
		u := sess.Get("User")

		if u == nil {
			panic("User is nil!")
		}

		user := u.(*models.User)

		pData := models.Page{ Title: "Startseite", User: user }

		tpl, _ := template.ParseFiles("html/frame_footer.html", "html/frame_header.html", "html/index.html")
		tpl.ExecuteTemplate(writer, "index.html", pData)
	})

	mux.HandleFunc("/login", handlers.HandleLogin)

	mux.Handle("/img", http.FileServer(http.Dir("img")))
	mux.Handle("/js", http.FileServer(http.Dir("js")))
	mux.Handle("/css", http.FileServer(http.Dir("css")))

	n.UseHandler(mux)
	n.Run(":8080")
}
