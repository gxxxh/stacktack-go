package worker

import amqp "github.com/rabbitmq/amqp091-go"

type HandleFunc func(deliveries <-chan amqp.Delivery, done chan error)

type ConsumerInterface interface {
	Run(handleFunc HandleFunc) error
	Shutdown(string) error
}
