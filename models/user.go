package models

import (
"fmt"
"regexp"
)

// User is a Struct that resembles a user containing
// UserId, Name, Username and HashedPassword
type User struct {
	Name               string
	Username, Password string
	Email              string
	HashedPassword     []byte
	Role               int
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")
