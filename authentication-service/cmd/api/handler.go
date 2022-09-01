package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type entry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// checks if server is running or not before sending request or getting response
func (app *Config) isAlive(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "im still alive",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

// authenticate the user by simple http and rest request
func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	req := request{}

	err := app.readJson(w, r, &req)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	usr, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	valid, err := usr.PasswordMatches(req.Password)
	if err != nil || !valid {
		app.errorJson(w, errors.New("invalid username or password"), http.StatusBadRequest)
		return
	}

	// log the user by logging service
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", usr.Email))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := response{
		Error:   false,
		Message: fmt.Sprintf("looged in %s", req.Email),
		Data:    usr,
	}

	app.writeJson(w, http.StatusAccepted, res)
}

// loggs for every person that authenticates and saves it in database
func (app *Config) logRequest(name, data string) error {
	entry := entry{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
