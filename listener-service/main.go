package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
	"log"
	"os"
	"time"
)

func main() {
	rabbitConn, err := connectRabbitmq()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("listening for rabbitmq messages")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{
		"log_INFO",
		"log_WARNING",
		"log_ERROR",
	})
	if err != nil {
		log.Println(err)
	}
}

func connectRabbitmq() (*amqp.Connection, error) {
	var (
		connection *amqp.Connection
		count      int
	)

	for {
		conn, err := amqp.Dial("ampq://guest:guest@rabbitmq")
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
