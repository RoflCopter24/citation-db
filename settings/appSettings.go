package settings

import "fmt"

type AppSettings struct {
	// DbServer is the hostname of the MongoDB server
	DbServer 	string
	// DbPort is the port of the MongoDB server
	DbPort		int
	// DbName is the name of the target DB on the MongoDB server
	DbName		string
	// DbUser is the username for accessing the MongoDB instance
	// "" if none
	DbUser		string
	// DbPass is the password required for MongoDB access
	// "" if none
	DbPass		string
	// Port of the CitationDB server
	Port		int
}

func (s *AppSettings) LffOrDefault() {
	s.Default()
}

func (s *AppSettings) Default() {
	// Default Mongo settings
	s.DbServer 	= "127.0.0.1"
	s.DbPort	= 27017
	s.DbName	= "citation"
	s.DbUser	= ""
	s.DbPass	= ""

	// App settings
	s.Port		= 8080
}

func (s *AppSettings) GenMgoConnStr() string {
	if s.DbUser == "" {
		return fmt.Sprintf("mongodb://%s:%d/%s", s.DbServer, s.DbPort, s.DbName)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", s.DbUser, s.DbPass, s.DbServer, s.DbPort, s.DbName)
}
