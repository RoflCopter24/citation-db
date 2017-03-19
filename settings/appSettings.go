package settings

import (
	"fmt"
	"os"
	"strconv"
)

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
	// Working directory
	WorkingDir	string
}

func (s *AppSettings) LffOrDefault() {
	s.Default()
}

func (s *AppSettings) Default() {
	// Default Mongo settings

	envServer := os.Getenv("MONGO_DB")
	if envServer == "" {
		envServer = "localhost"
	}

	envPort := os.Getenv("MONGO_DB_PORT")
	if envPort == "" {
		envPort = "27017"
	}

	s.DbServer 	= envServer
	s.DbPort,_	= strconv.Atoi(envPort)
	s.DbName	= "citation"
	s.DbUser	= ""
	s.DbPass	= ""

	// App settings
	s.WorkingDir	= os.Getenv("WORKINGDIR")
	s.Port		= 8080
}

func (s *AppSettings) GenMgoConnStr() string {
	if s.DbUser == "" {
		return fmt.Sprintf("mongodb://%s:%d/%s", s.DbServer, s.DbPort, s.DbName)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", s.DbUser, s.DbPass, s.DbServer, s.DbPort, s.DbName)
}
