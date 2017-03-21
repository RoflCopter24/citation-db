package handlers

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/RoflCopter24/citation-db/models"
	"html/template"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/RoflCopter24/negroni-sessions"
	"log"
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	tpl.ExecuteTemplate(w, "login.html", data)
}

func handleLoginPOST(w http.ResponseWriter, r *http.Request) {
	data := models.PageLogin{ Title: "Login", Success: false}

	err1 := r.ParseForm()
	if err1 != nil {
		log.Fatal(err1)
		data.Error= "Form data was invalid!"
		renderLoginTpl(data,w)
		return
	}

	userName := r.PostFormValue("username")
	userPass := r.PostFormValue("userpass")

	// You can access the mgo db object from the request object.
	// The db object is stored in key `db`.
	dbI := context.Get(r, "db")
	log.Println(context.Get(r, "mgoSession"))
	if dbI == nil {
		log.Fatal("Database Object is nil!")
		data.Error = "Database connection failed"
		renderLoginTpl(data, w)
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
	renderLoginTpl(data, w)
}

func renderLoginTpl(data models.PageLogin, w http.ResponseWriter) {
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
