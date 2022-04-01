package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"starbuy/authorization"
	"starbuy/model"
	"starbuy/repository"
	"starbuy/responses"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAddresses(w http.ResponseWriter, r *http.Request) {
	queried := mux.Vars(r)["user"]
	user, err := authorization.ExtractUser(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Token inválido"))
		return
	}

	if user != queried {
		responses.Error(w, http.StatusUnauthorized, errors.New("Não autorizado"))
		return
	}

	var addresses []model.RawAddress
	if err := repository.DownloadAddresses(user, &addresses); err != nil {
		if err == sql.ErrNoRows {
			responses.Error(w, http.StatusNoContent, errors.New("Nenhum endereço encontrado"))
			return
		}
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, addresses)
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	queried := mux.Vars(r)["user"]
	id := mux.Vars(r)["id"]

	user, err := authorization.ExtractUser(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Token inválido"))
		return
	}

	if user != queried {
		responses.Error(w, http.StatusUnauthorized, errors.New("Não autorizado"))
		return
	}

	var address model.Address
	if err := repository.DownloadAddress(id, &address); err != nil {
		if err == sql.ErrNoRows {
			responses.Error(w, http.StatusNoContent, errors.New("Nenhum endereço encontrado"))
			return
		}
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, address)
}

func PostAddress(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	user, err := authorization.ExtractUser(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Token inválido"))
		return
	}

	type Request struct {
		CEP        string `json:"cep,omitempty"`
		Number     int    `json:"number,omitempty"`
		Complement string `json:"complement,omitempty"`
	}

	var req Request
	if err = json.Unmarshal(body, &req); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	//TODO: Usar alguma API para verificar se o CEP bate com algum existente

	address := model.RawAddress{
		Identifier: strings.Replace(uuid.New().String(), "-", "", 4),
		Holder:     user,
		CEP:        req.CEP,
		Number:     req.Number,
		Complement: req.Complement,
	}

	responses.JSON(w, http.StatusOK, address)
}
