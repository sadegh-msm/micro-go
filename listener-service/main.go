package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
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
	log.Println("connected to rabbitmq")

}

func connectRabbitmq() (*amqp.Connection, error) {
	var (
		connection *amqp.Connection
		count      int
	)

	for {
		conn, err := amqp.Dial("ampq://guest:guest@localhost")
		if err != nil {
			fmt.Println("can not connect to rabbitmq")
			count++
		} else {
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
