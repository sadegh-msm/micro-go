package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func ResolveURL(c echo.Context) error {
	url := c.Param("url")

	r := CreateClients(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			log.Fatal("error closing the database connection")
		}
	}(r)

	value, err := r.Get(Context, url).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "short URL not found in database"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "cannot connect to database"})
	}

	rInr := CreateClients(1)
	defer func(rInr *redis.Client) {
		err := rInr.Close()
		if err != nil {
			log.Fatal("error closing the database connection")
		}
	}(rInr)

	_ = rInr.Incr(Context, "counter")

	return c.Redirect(301, value)
}
