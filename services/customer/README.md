# Customer Service

This microservice is responsible for managing user profiles, settings, and beneficiaries.

## Description

The Customer Service acts as the source of truth for user-related data after the initial onboarding. It consumes `user.created` events to create corresponding customer records in the Anchor BaaS. It also provides APIs for fetching user profiles, managing settings (like default receiving accounts), and handling the lifecycle of beneficiaries (external bank accounts).

## Endpoints

- `GET /users/me`: Fetches the profile of the authenticated user.
- `GET /users/{username}`: Fetches a public user profile.
- `POST /beneficiaries`: Adds a new external bank account for the user.
- `GET /beneficiaries`: Lists a user's saved beneficiaries.

## Dependencies

- Supabase (PostgreSQL)
- RabbitMQ
- Anchor API