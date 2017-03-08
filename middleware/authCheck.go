package middleware

import (
	"github.com/urfave/negroni"
	"net/http"
	"github.com/gorilla/context"
	"github.com/RoflCopter24/citation-db/models"
)

type AuthChecker struct {

}

func (ac *AuthChecker) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		granted := false

		context.Set(request, "AccessGranted",  granted)

		if !granted {
			context.Set(request, "TargetUrl", request.RequestURI)
		}


		context.Set(request, "User", models.User{})
		next(writer, request)
	}
}
