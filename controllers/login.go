package controllers

import (
	"database/sql"
	"net/http"
	"starbuy/authorization"
	"starbuy/database"
	"starbuy/model"
	"starbuy/repository"
	"starbuy/security"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func Auth(c *gin.Context) error {
	db := database.GrabDB()
	login := Login{}

	if err := c.BindJSON(&login); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad request"})
		return nil
	}

	recorded := Login{}
	if err := db.Get(&recorded, "SELECT * FROM login WHERE username=$1", login.Username); err != nil {
		if err == sql.ErrNoRows {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": "not found"})
			return nil
		}
		return err
	}

	var user model.User
	if err := repository.DownloadUser(login.Username, &user); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": "not found"})
		return nil
	}

	if err := security.ComparePassword(recorded.Password, login.Password); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "wrong password"})
		return nil
	}

	token := authorization.GenerateToken(login.Username)

	type Response struct {
		User  model.User `json:"user"`
		Token string     `json:"jwt"`
	}
	c.JSON(http.StatusOK, Response{User: user, Token: token})

	return nil
}
