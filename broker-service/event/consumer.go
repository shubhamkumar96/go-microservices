package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Define type that will be used for receiving from queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Define type that is used to pushing events to queue
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Define a method to consume messages
func (consumer *Consumer) Listen(topics []string) error {
	// Get a channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Get a random queue, and use it
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Bind the channel with each of the topics
	for _, topic := range topics {
		err := ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false, // no-wait ?
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}

	// Now we will consume messages
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Consume the messages from RabbitMQ, until the application closes.
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] : [logs_topic, %s]\n", q.Name)
	<-forever // Current go-routine will be blocked here forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	// Here we can add as many cases we want, as long as we write logic for its operation.

	default:
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logEvent(payload Payload) error {
	// create json that we will send to the logger-service
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// create a http request
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	// call the logger-service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// get the correct status-code
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
