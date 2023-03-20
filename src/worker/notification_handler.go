package worker

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// This is an example handler, which constains the essential parts of the handler function.
// 1. user need to define cleanup function and use the done channel to wakeup Shutdown function
// 2. if the AutoAck is false, user need to call d.Ack to acknowledge every deliviery
func SimpleHandler(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()
	for d := range deliveries {
		log.Println(d.Body)
		d.Ack(false)
	}
}
