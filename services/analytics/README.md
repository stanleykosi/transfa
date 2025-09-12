# Analytics Service

This microservice is responsible for processing transaction data and providing financial insights to users.

## Description

The Analytics Service operates asynchronously to avoid impacting the performance of core transactions. It consumes `transaction.completed` events from a message broker. Upon receiving an event, it updates pre-aggregated data tables that summarize user spending, cash flow, and other metrics. It then exposes endpoints for the iOS client to fetch this processed data quickly for display in the Analytics tab.

## Endpoints

- `GET /analytics/cash-flow`: Provides data for "Money In vs. Money Out" charts.
- `GET /analytics/spending-summary`: Provides data categorized by spending type.

## Dependencies

- Supabase (PostgreSQL) - for its own aggregate tables
- RabbitMQ