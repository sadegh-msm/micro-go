package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"time"
)

const port = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connectRabbitmq()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}
	log.Printf("starting broker service on port %s \n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// starting server
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectRabbitmq() (*amqp.Connection, error) {
	var (
		connection *amqp.Connection
		count      int
	)

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("can not connect to rabbitmq")
			count++
		} else {
			log.Println("connected to rabbitmq")
			connection = conn
			break
		}

		if count > 10 {
			fmt.Println(err)
			return nil, err
		}

		time.Sleep(2 * time.Second)
		log.Println("restarting connection")
		continue
	}

	return connection, nil
}
