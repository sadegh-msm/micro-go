package main

import (
	"logger-service/data"
	"net/http"
)

type request struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) isAlive(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "still alive",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

func (app *Config) writeLog(w http.ResponseWriter, r *http.Request) {
	var req request

	app.readJson(w, r, &req)

	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := response{
		Error:   false,
		Message: "logged data",
	}
	app.writeJson(w, http.StatusAccepted, res)
}
