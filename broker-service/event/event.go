package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// creating function for less typing while using them

func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
}

func declareRandomQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare("", false, false, true, false, nil)
}
