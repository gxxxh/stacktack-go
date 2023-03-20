package worker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

// const uri = "amqp://%s:%s@%s:%s/"
const uri = "amqp://%s:%s@%s:%s/"

var (
	ErrLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	Log    = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
)

func init() {

}

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	done       chan error //关闭标志
	Deliveries <-chan amqp.Delivery
}

func NewConsumer(rabbitMQConnectConfig *RabbitMQConnectConfig) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		done:    make(chan error),
	}
	amqpURI := fmt.Sprintf(uri, rabbitMQConnectConfig.RabbitUserID, rabbitMQConnectConfig.RabbitPassword, rabbitMQConnectConfig.RabbitHost, rabbitMQConnectConfig.RabbitPort)
	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(rabbitMQConnectConfig.Name)
	Log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}
	Log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	//连接被动关闭时运行
	go func() {
		Log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()
	return c, nil
}

func (c *Consumer) createExchange(exchangeConfig *ExchnageConfig) error {
	Log.Printf("Declaring exchange %s", exchangeConfig.Exchange)
	if err := c.channel.ExchangeDeclare(
		exchangeConfig.Exchange,     // name of the exchange
		exchangeConfig.ExchangeType, // type
		exchangeConfig.Durable,      // durable
		exchangeConfig.AutoDelete,   // delete when complete
		exchangeConfig.Internal,     // internal
		exchangeConfig.NoWait,       // noWait
		exchangeConfig.Args,         // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}

func (c *Consumer) createQueue(queueConfig *QueueConfig) error {
	Log.Printf("Declaring Queue %s", queueConfig.QueueName)
	_, err := c.channel.QueueDeclare(
		queueConfig.QueueName,  // name of the queue
		queueConfig.Durable,    // durable
		queueConfig.AutoDelete, // delete when unused
		queueConfig.Exclusive,  // exclusive
		queueConfig.NoWait,     // noWait
		queueConfig.Args,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	return nil
}

func (c *Consumer) Bind(exchangeConfig *ExchnageConfig, queueConfig *QueueConfig) error {
	var err error
	if err = c.createExchange(exchangeConfig); err != nil {
		return err
	}
	if err = c.createQueue(queueConfig); err != nil {
		return err
	}

	Log.Printf("declared Queue (%q), binding to Exchange (key %q)",
		queueConfig.QueueName, queueConfig.BindingKey)
	if err = c.channel.QueueBind(
		queueConfig.QueueName,   // name of the queue
		queueConfig.BindingKey,  // bindingKey
		exchangeConfig.Exchange, // sourceExchange
		queueConfig.NoWait,      // noWait
		queueConfig.Args,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}
	return nil
}
func (c *Consumer) StartConsume(exchangeConfig *ExchnageConfig, queueConfig *QueueConfig, consumerTag string) error {
	var err error
	Log.Printf("starting Consume (consumer tag %s ) Exhange %s on queue %s", consumerTag, exchangeConfig.Exchange, queueConfig.QueueName)

	c.Deliveries, err = c.channel.Consume(
		queueConfig.QueueName, // name
		consumerTag,           // consumerTag,
		false,                 // autoAck
		false,                 // exclusive
		false,                 // noLocal
		false,                 // noWait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}
	return nil
}

// 用户主动关闭Consumer
func (c *Consumer) Shutdown(consumerTag string) error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(consumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer Log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return c.WaitingForCleanUp()
}

func (c *Consumer) CleanUp() {
	Log.Printf("handle: deliveries channel closed")
	c.done <- nil
}

func (c *Consumer) WaitingForCleanUp() error {
	return <-c.done
}
