package main

import (
	"net/http"
)

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
