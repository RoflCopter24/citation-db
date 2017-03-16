package handlers

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/RoflCopter24/citation-db/models"
	"html/template"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/goincremental/negroni-sessions"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		handleLoginGET(w, r)
	} else if r.Method == "POST" {
		handleLoginPOST(w, r)
	}
}

func handleLoginGET(w http.ResponseWriter, r *http.Request) {
	data := models.PageLogin{ Title: "Login", Success: false }

	tpl, err := template.ParseFiles("html/frame_footer.html", "html/frame_header.html", "html/login.html")
	if err != nil {
		panic(err)
	}
	tpl.ExecuteTemplate(w, "login.html", data)
}

func handleLoginPOST(w http.ResponseWriter, r *http.Request) {
	data := models.PageLogin{ Title: "Login", Success: false}

	err1 := r.ParseForm()
	if err1 != nil {
		panic(err1)
	}

	userName := r.PostFormValue("username")
	userPass := r.PostFormValue("userpass")

	// You can access the mgo db object from the request object.
	// The db object is stored in key `db`.
	dbI := context.Get(r, "db")
	if dbI == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	db := dbI.(*mgo.Database)
	// Now lets perform a count query using mgo db object.

	user := models.User{}
	err := db.C("users").Find(bson.M{ "username": userName }).One(&user)

	if err != nil {
		data.Error = err.Error()
	} else {
		err2 := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(userPass))
		if err2 != nil {
			data.Error = err2.Error()
		} else {
			data.Success = true
			data.User = &user
			session := sessions.GetSession(r)
			session.Set("Username", userName)
			//session.Set("User", &user)
			context.Set(r, "User", userName)
			//context.Set(r, "User", &user)

			tUrl := session.Get("TargetUrl")
			if tUrl == nil {
			tUrl = "/"
			}
			data.TargetUrl = tUrl.(string)
		}
	}

	// You can access the mgo session object from the request object.
	// The session object is stored in key `mgoSession`.
	// mgoSession := context.Get(r, "mgoSession").(*mgo.Session)
	// count2, _ := mgoSession.DB("citation").C("users").Find(bson.M{ "Username": "Maik" }).Count()

	tpl, err := template.ParseFiles("html/frame_footer.html", "html/frame_header.html", "html/login.html")
	if err != nil {
		panic(err)
	}
	tpl.ExecuteTemplate(w, "login.html", data)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Clear()
	context.Clear(r)

	data := models.Page{ Title: "Abmeldung"}

	tpl, err := template.ParseFiles("html/frame_footer.html", "html/frame_header.html", "html/logout.html")
	if err != nil {
		panic(err)
	}
	tpl.ExecuteTemplate(w, "logout.html", data)
}
