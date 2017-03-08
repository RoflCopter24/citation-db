package middleware

import (
	"github.com/urfave/negroni"
	"net/http"
	"github.com/gorilla/context"
	"github.com/RoflCopter24/citation-db/models"
	//"gopkg.in/mgo.v2"
	"github.com/goincremental/negroni-sessions"
	"strings"
)

type AuthChecker struct {
	UsersDB		string
	UsersCollection	string
	WhiteList	[]string
}

func (ac *AuthChecker) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {

		// Ignore whitelisted urls
		for i:=range ac.WhiteList {
			if strings.HasPrefix(request.RequestURI, ac.WhiteList[i]) {
				next(writer, request)
				return
			}
		}

		// Session is stored in memory on the server, so we can put the
		// whole User object there and retrieve it
		session := sessions.GetSession(request)
		u := session.Get("User")

		if u == nil {
			context.Set(request, "TargetUrl", request.RequestURI)
			session.Set("TargetUrl", request.RequestURI)
			http.Redirect(writer, request, "/login", 401)
			return
		}

		user := u.(*models.User)
		context.Set(request, "User", user)
		next(writer, request)
	}
}
