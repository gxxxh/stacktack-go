package worker

import amqp "github.com/rabbitmq/amqp091-go"

func NovaNotificationInfoHandler(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()
	for d := range deliveries {
		if true {
			Log.Printf(
				"got %dB delivery: [%v] %q",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			)
		}
		d.Ack(false)
	}

}

func NovaNotificationErrorHandler(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()
	for d := range deliveries {
		Log.Printf(
			"got %dB delivery: [%v] %s",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
}
