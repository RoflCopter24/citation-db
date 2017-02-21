package controllers

import (
	"github.com/RoflCopter24/citation-db/app/routes"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	authenticated bool
}

func (c App) Index() revel.Result {
	if c.Session["user"] != "Mike" {
		return c.Redirect(routes.App.Login())
	}
	return c.RenderTemplate("App/Index.html")
}

func (c App) Login() revel.Result {
	userName := c.Params.Get("inputUser")
	password := c.Params.Get("inputPassword")
	revel.WARN.Println("User: " + userName)
	revel.WARN.Println("Pass: " + password)

	userName = c.Params.Form.Get("inputUser")
	password = c.Params.Form.Get("inputPassword")
	revel.WARN.Println("User: " + userName)
	revel.WARN.Println("Pass: " + password)

	if userName == "Mike" {
		if password == "test123" {
			c.authenticated = true
			c.Session["user"] = "Mike"
			c.Flash.Success("Welcome, " + "Mike")
			return c.Redirect(routes.App.Index())
		}
		return c.Forbidden("Access denied")
	}
	return c.Forbidden("Wrong credentials")
}

func (c App) LoginPage() revel.Result {

	return c.RenderTemplate("App/login.html")
}
