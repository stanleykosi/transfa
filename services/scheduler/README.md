# Scheduler Service

This microservice is a background worker responsible for running scheduled, time-based tasks (cron jobs).

## Description

The Scheduler Service does not typically expose a public API. Instead, it runs a set of predefined jobs on a recurring schedule. Its key responsibilities include:
- Debiting user wallets for monthly subscription fees.
- Resetting monthly free transfer limits for non-subscribed users.
- Processing expired Money Drops and returning remaining funds to the creator.

## Endpoints

This service has no public endpoints other than a `/health` check for monitoring.

## Dependencies

- Supabase (PostgreSQL)
- Transaction Service (to initiate transfers)
- Subscription Service (to get subscription data)