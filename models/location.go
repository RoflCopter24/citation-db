package models

type Location struct {
	Name 	string `bson:"name"`
	Notes	string `bson:"notes"`
}
