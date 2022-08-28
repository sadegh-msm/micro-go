package main

import (
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

func (app *Config) Broker(c echo.Context) error {
	res := response{
		Error:   false,
		Message: "Broker is active",
	}

	return c.JSON(http.StatusAccepted, res)
}
