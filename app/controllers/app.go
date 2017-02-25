package controllers

import (
	"github.com/RoflCopter24/citation-db/app/models"
	"github.com/RoflCopter24/citation-db/app/routes"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type App struct {
	MongoController
}

func (c App) AddUser() revel.Result {

	if user := c.IsUserLoggedIn(); user != nil {
		c.RenderArgs["user"] = user
	}

	return nil
}

func (c App) IsUserLoggedIn() *models.User {

	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}

	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c App) getUser(username string) *models.User {
	result := models.User{}
	coll := c.mgoSess.DB("citation").C("users")
	revel.TRACE.Printf("Username is %s\n", username)
	revel.TRACE.Print(coll.Name)
	err := coll.Find(bson.M{"Username": username}).One(&result)

	if err != nil {
		revel.ERROR.Fatal(err)
	}

	return &result
}

func (c App) Index() revel.Result {
	if user := c.IsUserLoggedIn(); user == nil {
		return c.Redirect(routes.App.LoginPage())
	}
	return c.RenderTemplate("App/Index.html")
}

func (c App) Login(remember bool) revel.Result {
	username := c.Params.Get("inputUser")
	password := c.Params.Get("inputPassword")
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.App.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.LoginPage())
}

func (c App) LoginPage() revel.Result {
	return c.RenderTemplate("App/login.html")
}

/* REGISTER USER */
func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.LoginPage())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	coll := s.DB("citation").C("users")
	err := coll.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}
