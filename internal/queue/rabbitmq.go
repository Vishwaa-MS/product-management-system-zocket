package queue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ represents the RabbitMQ connection and channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
}

// NewRabbitMQ initializes a new RabbitMQ instance
func NewRabbitMQ(host string, port int, queueName string) *RabbitMQ {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%d/", host, port))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		QueueName:  queueName,
	}
}

// PublishMessage publishes a message to the queue
func (r *RabbitMQ) PublishMessage(message string) error {
	return r.Channel.Publish(
		"",            // exchange
		r.QueueName,   // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// ConsumeMessages starts consuming messages from the queue
func (r *RabbitMQ) ConsumeMessages() <-chan string {
	msgs, err := r.Channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Convert messages to a string channel
	out := make(chan string)
	go func() {
		defer close(out)
		for d := range msgs {
			out <- string(d.Body)
		}
	}()
	return out
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.Connection.Close()
}
