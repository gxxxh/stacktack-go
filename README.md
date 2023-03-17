# stacktack-go

This library aims to monitor notifications generated by openstack.

## Usage

### Configs

User need to define a config file, here is an example.

```json
{
  "RabbitMQConnectConfig": {
    "Name": "NovaInfoConsumer",
    "RabbitHost": "133.133.135.136",
    "RabbitPort": "5672",
    "RabbitUserID": "guest",
    "RabbitPassword": "guest",
    "RabbitVirtualHost": "/"
  },
  "ExchangeConfig": {
    "Exchange": "nova",
    "ExchangeType": "topic",
    "Durable": false,
    "AutoDelete": false,
    "Internal": false,
    "NoWait": false
  },
  "QueueConfig": {
    "QueueName": "notifications.info",
    "BindingKey": "notifications.info",
    "Durable": false,
    "AutoDelete": false,
    "Exclusive": false,
    "NoWait": false
  },
  "ConsumerTag": "NovaInfoConsumer",
  "AutoAck": false
}
```

User can use the following command to get the configs

#### RabbitMQConnectConfig

```shell
cat /etc/rabbitmq/rabbitmq.config
cat /etc/rabbitmq/rabbitmq-env.conf
# 查看使用的端口
lsof -i | grep beam
```

#### Binding Config

```shell
rabbitmqctl list_bindings | grep nova
```

#### Exchange Config

```shell
rabbitmqctl list_exchanges name type durable auto_delete internal arguments | grep nova
```

#### Queue Config

```shell
rabbitmqctl list_queues name durable auto_delete arguments policy  exclusive exclusive_consumer_pid exclusive_consumer_tag| grep notific
```

#### firewalls

To make this library working, user need to open the port by firewalld or iptables

```shell
```

### Use

Here is an example to Consume the notifications from openstack.

```shell
func TestNovaConsumerRun(t *testing.T) {
	config := worker.NewConfigByJson("E:\\gopath\\src\\stacktack-go\\test\\worker\\config.json")

	novaConsumer, err := worker.NewNovaConsumer(config)
	if err != nil {
		t.Error(err)
		return
	}
	novaConsumer.Run(worker.NovaNotificationInfoHandler)
}
```

Notice that user needs to define a notificationHandler function to process the notifications, Here is an example handler
function.

```shell
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
```