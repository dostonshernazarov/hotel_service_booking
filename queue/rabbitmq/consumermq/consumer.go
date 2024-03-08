package consumermq

import (
	"github.com/streadway/amqp"
)

type RabbitMQConsumer interface {
	ConsumerMassages(handler func(message []byte)) error
	Close() error
}

type RabbitMQConsumerImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQConsumer(amqpURI, queueName string) (RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumerImpl{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (r *RabbitMQConsumerImpl) ConsumerMassages(handler func(message []byte)) error {
	msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

func (r *RabbitMQConsumerImpl) Close() error {
	return r.conn.Close()
}
