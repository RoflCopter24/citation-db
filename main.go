package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
	"github.com/gorilla/context"
	mgosession "github.com/joeljames/nigroni-mgo-session"
	mgo "gopkg.in/mgo.v2"
	"github.com/RoflCopter24/citation-db/settings"
	"github.com/RoflCopter24/citation-db/middleware"
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

	// Setup MongoDb connection stuff
	setupMgo(n, &appSettings)

	ac := middleware.AuthChecker{}
	n.Use(ac.Middleware())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// You can access the mgo db object from the request object.
		// The db object is stored in key `db`.
		db := context.Get(request, "db").(*mgo.Database)
		// Now lets perform a count query using mgo db object.
		count, _ := db.C("users").Find(nil).Count()
		fmt.Fprintf(writer, "Determining the count in the collection using the db object. \n\n")
		fmt.Fprintf(writer, "Total number of object in the mongo database: %d  \n\n", count)

		// You can access the mgo session object from the request object.
		// The session object is stored in key `mgoSession`.
		mgoSession := context.Get(request, "mgoSession").(*mgo.Session)
		count2, _ := mgoSession.DB(appSettings.DbName).C("users").Find(nil).Count()
		fmt.Fprintf(writer, "Determining the count in the collection using the session object. \n\n")
		fmt.Fprintf(writer, "Total number of object in the mongo database: %d  \n\n", count2)

	})

	n.UseHandler(mux)
	n.Run(":8080")
}
