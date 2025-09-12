# Account Service

This microservice is responsible for creating and managing user accounts (wallets) within the BaaS provider.

## Description

The Account Service listens for `customer.verified` events. Upon receiving an event, it communicates with the Anchor API to provision a `DepositAccount` for the user, which serves as their primary in-app wallet. It then stores the account details and ID in the Supabase database. It also handles the creation of special-purpose accounts, like the persistent Money Drop wallet.

## Endpoints

This service is primarily event-driven and may not expose public HTTP endpoints initially, other than for internal health checks.

## Dependencies

- Supabase (PostgreSQL)
- RabbitMQ
- Anchor API