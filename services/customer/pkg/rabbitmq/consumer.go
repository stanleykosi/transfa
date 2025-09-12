/**
 * @description
 * This package provides a reusable RabbitMQ consumer. It handles the complexities of
 * connecting to RabbitMQ, setting up a channel, declaring queues and exchanges,
 * and consuming messages in a resilient way.
 *
 * Key features:
 * - Establishes and maintains a connection to RabbitMQ.
 * - Declares necessary topology (exchange, queue, binding).
 * - Starts a message consumer in a separate goroutine.
 * - Provides a clean shutdown mechanism.
 *
 * @dependencies
 * - "context", "fmt", "log"
 * - "github.com/rabbitmq/amqp091-go": The official Go client for RabbitMQ.
 */
package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// Consumer holds the necessary components for a RabbitMQ consumer.
type Consumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// MessageHandler is a function type that processes a single RabbitMQ message.
type MessageHandler func(ctx context.Context, msg amqp091.Delivery) error

// NewConsumer creates and returns a new RabbitMQ consumer.
func NewConsumer(amqpURL string) (*Consumer, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
	}, nil
}

// StartConsumer sets up the RabbitMQ topology and begins consuming messages.
func (c *Consumer) StartConsumer(ctx context.Context, exchange, queueName, routingKey, consumerTag string, handler MessageHandler) error {
	// Declare the exchange
	err := c.channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare the queue
	q, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind the queue to the exchange
	err = c.channel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchange,     // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	// Start consuming messages
	msgs, err := c.channel.Consume(
		q.Name,       // queue
		consumerTag,  // consumer
		false,        // auto-ack is false, we will manually acknowledge
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Run message processing in a separate goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer shutting down...")
				c.Close()
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed. Exiting consumer goroutine.")
					return
				}
				log.Printf("Received a message: %s", msg.Body)
				if err := handler(ctx, msg); err != nil {
					log.Printf("Error processing message: %v. Nacking.", err)
					// Negative Acknowledge the message, requeue for another attempt
					msg.Nack(false, true)
				} else {
					// Acknowledge the message was processed successfully
					msg.Ack(false)
					log.Printf("Message processed and acknowledged.")
				}
			}
		}
	}()

	log.Printf("Consumer started. Waiting for messages on queue '%s'", queueName)
	return nil
}

// Close gracefully closes the channel and connection.
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}