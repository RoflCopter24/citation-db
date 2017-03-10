package models

import (
	"math/rand"
)

type Page struct {
	Title 		string
	User		*User
	CheckStr 	string
}

func (p *Page) GenCheckStr() {
	p.CheckStr = string(rand.Int())
}
