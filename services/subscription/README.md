# Subscription Service

This microservice manages user subscription plans, status, and entitlements.

## Description

The Subscription Service is the source of truth for whether a user is on the free or paid tier. It provides APIs for other services to query a user's subscription status and their usage of metered features (like free monthly external transfers). It also exposes endpoints for the client to allow users to upgrade, downgrade, or cancel their subscriptions.

## Endpoints

- `GET /subscriptions/me`: Gets the authenticated user's subscription status.
- `POST /subscriptions/subscribe`: Upgrades a user to the paid tier.
- `POST /subscriptions/cancel`: Cancels a user's subscription renewal.

## Dependencies

- Supabase (PostgreSQL)
- Scheduler Service (for billing triggers)