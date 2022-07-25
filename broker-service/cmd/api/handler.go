package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func isAlive(c echo.Context) error {
	return c.JSON(http.StatusOK, "still alive")
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	res := response{
		Error:   false,
		Message: "Broker is active",
	}

	out, _ := json.MarshalIndent(res, "", "\t")
	
}