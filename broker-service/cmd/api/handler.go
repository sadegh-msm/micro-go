package main

import (
	"bytes"
	"encoding/json"
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
		app.Authenticate(w, req.Auth)
	default:
		app.errorJson(w, errors.New("invalid action"), http.StatusBadRequest)
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, data AuthPayload) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("user Unauthorized"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error when authorizing"))
		return
	}

	var responseFromService response

	err = json.NewDecoder(res.Body).Decode(&responseFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if responseFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	payload := response{
		Error:   false,
		Message: "auth completed",
		Data:    responseFromService,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}
