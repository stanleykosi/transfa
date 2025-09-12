/**
 * @description
 * This package provides a concrete implementation of the Publisher interface for RabbitMQ.
 * It handles the connection, channel management, and message publishing logic,
 * abstracting away the complexities of the AMQP protocol from the core application.
 *
 * @dependencies
 * - "context": For context-aware operations.
 * - "encoding/json": To serialize event payloads.
 * - "log": For logging errors.
 * - "github.com/rabbitmq/amqp091-go": The RabbitMQ client library.
 */
package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher implements the app.Publisher interface for RabbitMQ.
type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewPublisher creates and returns a new RabbitMQ publisher.
// It establishes a connection and opens a channel.
func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Publisher{
		conn: conn,
		ch:   ch,
	}, nil
}

// Publish sends a message to a specific exchange with a given routing key.
// It ensures the exchange exists before publishing.
func (p *Publisher) Publish(ctx context.Context, body []byte, exchange, routingKey string) error {
	// Declare a topic exchange to ensure it exists.
	// Topic exchanges are durable and auto-deleted is false.
	err := p.ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Close gracefully closes the RabbitMQ channel and connection.
func (p *Publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	log.Println("RabbitMQ publisher connection closed.")
}