# Notification Service

This microservice handles all incoming webhooks and outgoing user notifications.

## Description

The Notification Service is the primary interface for external, asynchronous events. It exposes a webhook endpoint to receive events from the Anchor BaaS, verifies their authenticity, and then publishes them as internal events onto RabbitMQ for other services to consume (e.g., `customer.identification.approved`). It also consumes internal events to send user-facing notifications, such as push notifications or emails.

## Endpoints

- `POST /webhooks/anchor`: Receives and processes webhooks from the Anchor BaaS.

## Dependencies

- RabbitMQ
- Anchor API (for webhook signature verification)
- Apple Push Notification Service (APNS) (or other push notification providers)