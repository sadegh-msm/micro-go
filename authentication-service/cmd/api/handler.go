package main

import (
	"errors"
	"fmt"
	"net/http"
)

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) isAlive(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "im still alive",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

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

	res := response{
		Error:   false,
		Message: fmt.Sprintf("looged in %s", req.Email),
		Data:    usr,
	}

	app.writeJson(w, http.StatusAccepted, res)
}
