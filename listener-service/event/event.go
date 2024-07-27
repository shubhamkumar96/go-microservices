package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type of the exchange
		true,         // is the exchange durable
		false,        // should auto-delete, if not in use
		false,        // is this used internally
		false,        // no-wait ?
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name of the queue, pick a name randomly
		false, // is the queue durable
		false, // should auto-delete, if not in use
		true,  // is this exclusive queue
		false, // no-wait ?
		nil,   // arguments
	)
}
