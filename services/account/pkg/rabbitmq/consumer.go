/**
 * @description
 * This file provides a generic RabbitMQ consumer client. It handles the boilerplate
 * of connecting to RabbitMQ, declaring exchanges and queues, binding them, and
 * consuming messages.
 *
 * @dependencies
 * - Standard library packages for context and logging.
 * - "github.com/rabbitmq/amqp091-go": The official RabbitMQ client library for Go.
 */
package rabbitmq

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// Consumer holds the necessary components for a RabbitMQ consumer.
type Consumer struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

// MessageHandler is a function type that processes a delivered message.
// It returns an error to indicate if the message processing failed.
type MessageHandler func(ctx context.Context, msg amqp091.Delivery) error

// NewConsumer creates and returns a new RabbitMQ consumer.
func NewConsumer(amqpURL string) (*Consumer, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{conn: conn, ch: ch}, nil
}

// StartConsumer declares the necessary topology (exchange, queue, binding)
// and starts consuming messages, passing them to the provided handler.
func (c *Consumer) StartConsumer(ctx context.Context, exchange, queueName, routingKey, consumerTag string, handler MessageHandler) error {
	// Declare a durable, topic-based exchange.
	err := c.ch.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Declare a durable queue.
	q, err := c.ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange with the routing key.
	err = c.ch.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Start consuming messages from the queue.
	msgs, err := c.ch.Consume(
		q.Name,
		consumerTag,
		false, // auto-ack is false, we will manually acknowledge.
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Run the message processing loop in a separate goroutine.
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer context cancelled, stopping message processing.")
				return
			case d, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed, stopping consumer.")
					return
				}
				err := handler(ctx, d)
				if err != nil {
					log.Printf("Error handling message: %v. Nacking message.", err)
					// Negative Acknowledge the message. 'requeue=false' sends it to a dead-letter queue if configured.
					d.Nack(false, false)
				} else {
					// Acknowledge the message was processed successfully.
					d.Ack(false)
				}
			}
		}
	}()

	log.Printf("Consumer started. Waiting for messages on queue '%s' with routing key '%s'.", queueName, routingKey)
	return nil
}

// Close gracefully closes the channel and connection to RabbitMQ.
func (c *Consumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}