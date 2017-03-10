package models

import (
	"math/rand"
	"strconv"
)

type Page struct {
	Title 		string
	User		*User
	CheckStr 	string
}

func (p *Page) GenCheckStr() {
	p.CheckStr = strconv.Itoa(rand.Int())
}
