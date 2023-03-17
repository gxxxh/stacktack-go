package worker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)

// This is a handler for queue notifications.info in nova, we use gjson to parse the delivery body.
// Notice that body is a raw json while the "olso.message" is a json string, and the user has to call
// gjson.Parse(nessageRawJson.str) to get the content of the json string
func NovaNotificationInfoHandler(deliveries <-chan amqp.Delivery, done chan error) {
	cleanup := func() {
		Log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()
	for d := range deliveries {
		bodyJson := gjson.ParseBytes(d.Body)
		messageJson := gjson.Parse(bodyJson.Get("oslo\\.message").Str)
		eventType := messageJson.Get("event_type")
		fmt.Println("event type is ", eventType)
		if strings.Contains(eventType.String(), "instance.") {
			payLoad := messageJson.Get("payload")
			// instance UUID's seem to hide in a lot of odd places.
			instanceID := ""
			if payLoad.Get("instance_id").Exists() {
				instanceID = payLoad.Get("instance_id").String()
			} else if payLoad.Get("instance_uuid").Exists() {
				instanceID = payLoad.Get("instacnce_uuid").String()
			} else if payLoad.Get("exception.kwargs.uuid").Exists() {
				instanceID = payLoad.Get("exception.kwargs.uuid").String()
			} else if payLoad.Get("instance.uuid").Exists() {
				instanceID = payLoad.Get("instance.uuid").String()
			}
			fmt.Println("Handling Event ", eventType.String(), ",instanceID is:", instanceID)
		}
		d.Ack(false)
	}

}

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
