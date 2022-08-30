package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/:url", ResolveURL)
	e.POST("/api/v1", ShortenURL)

	log.Fatal(e.Start(":80"))
}
