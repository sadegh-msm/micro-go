package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct{}

func main() {
	app := Config{}
	log.Printf("starting broker service on port %s \n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// starting server
	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
