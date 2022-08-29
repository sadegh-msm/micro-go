package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/:url", ResolveURL)
	e.POST("/api/v1", ShortenURL)

	log.Fatal(e.Start(os.Getenv("APP_PORT")))
}
