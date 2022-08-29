package main

import (
	"errors"
	"net/http"
)

type BrokerRequest struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) isAlive(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "im still alive",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	res := response{
		Error:   false,
		Message: "Broker is active",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

func (app *Config) HandleAll(w http.ResponseWriter, r *http.Request) {
	var req BrokerRequest

	err := app.readJson(w, r, &req)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "auth":

	default:
		app.errorJson(w, errors.New("invalid action"), http.StatusBadRequest)
	}
}
