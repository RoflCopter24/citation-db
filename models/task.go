package models

import "time"

type Task struct {
	Name 			string		`bson:"name"`
	Content			string		`bson:"content"`
	ExpirationDate		time.Time	`bson:"expirationDate"`
	ShowNotification	bool		`bson:"showNotification"`
	UserId			string 		`bson:"userId"`
}
