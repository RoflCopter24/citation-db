package controllers

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/RoflCopter24/citation-db/app/models"
	"github.com/revel/revel"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	s *mgo.Session
)

func InitDB() {
	revel.INFO.Println("Connecting to MongoDB on localhost")
	s2, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	s = s2
	defer s.Close()
	revel.INFO.Println("Connection successful.")

	// Optional. Switch the session to a monotonic behavior.
	s.SetMode(mgo.Monotonic, true)
	revel.TRACE.Println("Mongo session mode set to 'monotonic'.")

	c := s.DB("citation").C("users")

	revel.TRACE.Println("Checking if user collection is populated...")
	// Check if testuser exists
	result := models.User{}
	err = c.Find(bson.M{"Username": "Mike"}).One(&result)

	//If not populate default database values
	if err != nil {

		revel.TRACE.Println("User collection is NOT populated. Doing that now...")

		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("test123"), bcrypt.DefaultCost)

		err = c.Insert(&models.User{Name: "Mike Muster", Username: "Mike", Email: "test@florian-vick.de", HashedPassword: bcryptPassword, Role: 2})

		if err != nil { //If this fails, there is already an entry
			panic(err)
		}
	}
	revel.TRACE.Println("Database setup complete.")
}

type MongoController struct {
	*revel.Controller
	mgoSess *mgo.Session
}

func (mc MongoController) Begin() revel.Result {

	if s != nil {
		mc.mgoSess = s.Clone()
		return nil
	}

	s, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	mc.mgoSess = s
	return nil

}

func (mc MongoController) Commit() revel.Result {

	if mc.mgoSess != nil {
		mc.mgoSess.Close()
		return nil
	}
	return nil
}
