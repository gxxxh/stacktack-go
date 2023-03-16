package worker

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

type RabbitMQConnectConfig struct {
	Name              string
	RabbitHost        string
	RabbitPort        string
	RabbitUserID      string
	RabbitPassword    string
	RabbitVirtualHost string
}

// parameter used in ExchangeDeclare
type ExchnageConfig struct {
	Exchange     string
	ExchangeType string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Args         amqp.Table
}

// parameter used in QueueDeclare
type QueueConfig struct {
	QueueName  string
	BindingKey string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ConsumerConfig struct {
	RabbitMQConnectConfig *RabbitMQConnectConfig
	ExchangeConfig        *ExchnageConfig
	QueueConfig           *QueueConfig
	ConsumerTag           string //AMQP consumer tag (should not be blank)
	AutoAck               bool   //enable message auto-ack
}

func NewConfig() *ConsumerConfig {
	return &ConsumerConfig{}
}

func NewConfigByJson(fileName string) *ConsumerConfig {
	c := &ConsumerConfig{}
	c.Load(fileName)
	return c
}

func (c *ConsumerConfig) Load(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
}
