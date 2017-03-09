package models

import (
"fmt"
"regexp"
)

// User is a Struct that resembles a user containing
// UserId, Name, Username and HashedPassword
type User struct {
	FirstName		string 	`bson:"firstname"`
	LastName               	string 	`bson:"lastname"`
	Username           	string 	`bson:"username"`
	Email              	string 	`bson:"email"`
	HashedPassword     	[]byte 	`bson:"hashedpassword"`
	Role               	int	`bson:"role"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func IsValidUsername(name string) bool {
	return userRegex.MatchString(name)
}
