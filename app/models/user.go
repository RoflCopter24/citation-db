package models

import (
	"fmt"
	"regexp"

	"github.com/revel/revel"
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

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{Max: 15},
		revel.MinSize{Min: 4},
		revel.Match{Regexp: userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{Max: 100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{Max: 15},
		revel.MinSize{Min: 5},
	)
}
