/**
 * @description
 * This package provides a simple, reusable RabbitMQ publisher client.
 *
 * It abstracts the logic for connecting to RabbitMQ, declaring exchanges,
 * and publishing messages. This allows different services to send events
 * without duplicating connection and publishing logic.
 *
 * Key features:
 * - Manages a persistent connection and channel to RabbitMQ.
 * - Provides a simple `Publish` method to send messages.
 * - Handles graceful connection closing.
 *
 * @dependencies
 * - "context": For context-aware publishing.
 * - "fmt": For error formatting.
 * - "github.com/rabbitmq/amqp091-go": The official RabbitMQ Go client.
 */
package rabbitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// Publisher holds the connection and channel for publishing messages.
type Publisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewPublisher creates and returns a new Publisher instance.
// It establishes a connection to the RabbitMQ server using the provided URL.
func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Publisher{
		conn:    conn,
		channel: ch,
	}, nil
}

// Publish sends a message to a specified exchange with a routing key.
// It ensures the exchange exists by declaring it as a topic exchange.
func (p *Publisher) Publish(ctx context.Context, body []byte, exchange, routingKey string) error {
	// Ensure the exchange exists. This is idempotent.
	err := p.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

// Close gracefully closes the channel and connection to RabbitMQ.
func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}