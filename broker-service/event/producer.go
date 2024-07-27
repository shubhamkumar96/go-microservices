package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Define type for pushing events to RabbitMQ
type Producer struct {
	connection *amqp.Connection
}

func NewProducer(conn *amqp.Connection) (Producer, error) {
	producer := Producer{
		connection: conn,
	}

	err := producer.setup()
	if err != nil {
		return Producer{}, err
	}

	return producer, nil
}

func (p *Producer) setup() error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (p *Producer) Push(event, severity string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	// Pushing the event to channel
	err = channel.Publish(
		"logs_topic", severity, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
