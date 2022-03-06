package controllers

import (
	"authentication-service/database"
	"authentication-service/model"
	"authentication-service/responses"
	"authentication-service/security"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var data IncomingUser
	if err = json.Unmarshal(body, &data); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	user := model.User{data.Username, data.Email, data.Name, data.Gender, data.Birthdate, data.Seller, data.ProfilePicture, data.City, data.Registration}
	if err := user.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db := database.GrabDB()

	testQueries := map[string]string{
		fmt.Sprintf("SELECT * FROM users WHERE username='%s'", data.Username): "Username já está em uso",
		fmt.Sprintf("SELECT * FROM users WHERE email='%s'", data.Email):       "Email já está em uso",
	}

	var found model.User
	for key, value := range testQueries {
		err := db.Get(&found, key)
		if (err != nil && err != sql.ErrNoRows) || err == nil {
			responses.Error(w, http.StatusBadRequest, errors.New(value))
			return
		}
	}

	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO users VALUES (:username,:email,:name,:gender,:registration,:birthdate,:seller)", &user)
	if err := tx.Commit(); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	crypt, err := security.Hash(data.Password)
	tx2 := db.MustBegin()
	tx2.MustExec("INSERT INTO login VALUES ($1,$2)", data.Username, string(crypt))
	if err := tx2.Commit(); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, nil)
}
