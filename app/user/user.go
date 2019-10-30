package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type User struct {
	Name  string
	Email string
}

type randomUser struct {
	Results []struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email string `json:"email"`
	} `json:"results"`
}

func CreateUserHandler(logger *logrus.Logger, store Store) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		u := generate()
		if err := store.Save(u); err != nil {
			logger.WithFields(logrus.Fields{"error": err}).Error("Error handling Create User")
			c.JSON(500, gin.H{"message": "cannot create new user",})
			return
		}
		logger.WithFields(logrus.Fields{"name": u.Name, "email": u.Email}).Info("New user has been created")
		c.JSON(200, gin.H{"message": "ok",})
	}

	return gin.HandlerFunc(fn)
}

func generate() *User {

	res, _ := http.Get("https://randomuser.me/api/")
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	ru := randomUser{}
	_ = json.Unmarshal(body, &ru)

	return &User{
		Name:  ru.Results[0].Name.First + " " + ru.Results[0].Name.Last,
		Email: ru.Results[0].Email,
	}
}
