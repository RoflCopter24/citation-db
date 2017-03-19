package nigronimgosession

import (
	"net/http"

	"github.com/urfave/negroni"
)

type Database struct {
	dba DatabaseAccessor
}

func NewDatabase(databaseAccessor DatabaseAccessor) *Database {
	return &Database{databaseAccessor}
}

func (d *Database) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		reqSession := d.dba.Clone()
		defer reqSession.Close()
		d.dba.Set(request, reqSession)
		next(writer, request)
	}
}
