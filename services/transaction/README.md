# Transaction Service

This microservice is the core engine for all money movement within the Transfa application.

## Description

The Transaction Service handles the logic for all types of financial transactions. Its responsibilities include:
- Executing Peer-to-Peer (P2P) transfers, including the intelligent routing logic based on the recipient's subscription status.
- Processing Self-Transfers (withdrawals) to external bank accounts.
- Managing the entire lifecycle of Money Drops, from creation and funding to claims and expiry.
- Handling the creation and fulfillment of Payment Requests.
- Communicating with the Anchor API to initiate `BookTransfer` and `NIPTransfer` operations.

## Endpoints

- `POST /transactions/p2p`: Initiates a P2P transfer.
- `POST /transactions/self-transfer`: Initiates a withdrawal.
- `POST /money-drops`: Creates a new Money Drop.
- `POST /money-drops/{id}/claim`: Claims a Money Drop.
- `POST /payment-requests`: Creates a new Payment Request.

## Dependencies

- Supabase (PostgreSQL)
- RabbitMQ
- Anchor API
- Customer Service (to fetch recipient data)
- Subscription Service (to check subscription status)