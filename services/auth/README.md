# Auth Service

This microservice is responsible for handling user authentication and onboarding processes.

## Description

The Auth Service orchestrates the initial user creation after a successful signup via Clerk. It validates the Clerk JWT, creates a user record in the Supabase database, and publishes a `user.created` event to the message broker to trigger subsequent onboarding steps like customer creation in the BaaS.

## Endpoints

- `POST /onboarding`: Creates a new user profile after Clerk signup.

## Dependencies

- Supabase (PostgreSQL)
- RabbitMQ
- Clerk (for JWT validation)