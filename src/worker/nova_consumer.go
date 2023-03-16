package worker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NovaConsumer struct {
	*ConsumerConfig
	*Consumer
}

func NewNovaConsumer(consumerConfig *ConsumerConfig) (*NovaConsumer, error) {
	var err error
	nc := &NovaConsumer{}
	nc.ConsumerConfig = consumerConfig
	nc.Consumer, err = NewConsumer(consumerConfig.RabbitMQConnectConfig)
	if err != nil {
		return nil, err
	}
	err = nc.Bind(nc.ExchangeConfig, nc.QueueConfig)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func (nc *NovaConsumer) Run(handleFunc HandleFunc) error {
	//连接被动关闭时运行
	go func() {
		Log.Printf("closing: %s", <-nc.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	Log.Printf("starting Consume (consumer tag %q) on queue %s", nc.ConsumerTag, nc.QueueConfig.QueueName)
	deliveries, err := nc.channel.Consume(
		nc.QueueConfig.QueueName, // name
		nc.ConsumerTag,           // consumerTag,
		false,                    // autoAck
		false,                    // exclusive
		false,                    // noLocal
		false,                    // noWait
		nil,                      // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}
	handleFunc(deliveries, nc.done)
	return <-nc.done
}
