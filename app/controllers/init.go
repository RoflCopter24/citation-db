package controllers

import "github.com/revel/revel"

func init() {
	revel.OnAppStart(InitDB)
	//revel.InterceptMethod((*MongoController).Begin, revel.BEFORE)
	//revel.InterceptMethod(App.AddUser, revel.BEFORE)
	//revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	//revel.InterceptMethod((*MongoController).Commit, revel.AFTER)
}
