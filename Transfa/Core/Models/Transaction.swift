/**
 * @description
 * This file defines the core data model for a financial transaction.
 *
 * The `Transaction` struct represents a single money movement event within the Transfa system,
 * corresponding to a record in the `transactions` table. It's designed to be versatile enough
 * to represent various transaction types, from P2P payments to subscription fees.
 *
 * Key features:
 * - Conforms to `Codable` and `Identifiable` for API integration and SwiftUI compatibility.
 * - Uses enums for `TransactionType` and `TransactionStatus` for robust, type-safe logic.
 * - Handles various participants with optional UUIDs (e.g., a wallet funding transaction has no sender).
 * - Implements `CodingKeys` to map between Swift's camelCase and the backend's snake_case.
 *
 * @dependencies
 * - Foundation: Provides `UUID` and `Date` types.
 */
import Foundation

/// Represents a single financial transaction within the Transfa application.
struct Transaction: Codable, Identifiable {
    /// The unique identifier for the transaction.
    let id: UUID
    
    /// The user who initiated the transaction. Optional for system-generated transactions like wallet funding.
    let senderUserID: UUID?
    
    /// The user who received the funds. Optional.
    let recipientUserID: UUID?
    
    /// The internal Transfa account from which funds were debited. Optional.
    let sourceAccountID: UUID?
    
    /// The internal Transfa account to which funds were credited. Optional for external transfers.
    let destinationAccountID: UUID?
    
    /// The external beneficiary account to which funds were sent. Optional for internal transfers.
    let destinationBeneficiaryID: UUID?
    
    /// The unique ID for the transfer from the BaaS provider (Anchor). Optional.
    let anchorTransferID: String?
    
    /// The type of the transaction.
    let type: TransactionType
    
    /// The principal amount of the transaction, in the smallest currency unit (kobo).
    let amount: Int64
    
    /// The fee charged for the transaction, in the smallest currency unit (kobo).
    let fee: Int64
    
    /// The current status of the transaction.
    let status: TransactionStatus
    
    /// A user-provided description for the transaction. Optional.
    let description: String?
    
    /// A system- or user-assigned category for analytics. Optional.
    let category: String?
    
    /// The timestamp when the transaction was created.
    let createdAt: Date
    
    /// The timestamp when the transaction was last updated.
    let updatedAt: Date
    
    /// Defines the possible types of transactions. Raw values match the database enum.
    enum TransactionType: String, Codable {
        case p2p
        case selfTransfer = "self_transfer"
        case moneyDropFunding = "money_drop_funding"
        case moneyDropClaim = "money_drop_claim"
        case subscriptionFee = "subscription_fee"
        case walletFunding = "wallet_funding"
    }

    /// Defines the possible states of a transaction. Raw values match the database enum.
    enum TransactionStatus: String, Codable {
        case pending
        case completed
        case failed
        case reversed
    }
    
    /// Maps Swift's camelCase properties to the database's snake_case JSON keys.
    enum CodingKeys: String, CodingKey {
        case id
        case senderUserID = "sender_user_id"
        case recipientUserID = "recipient_user_id"
        case sourceAccountID = "source_account_id"
        case destinationAccountID = "destination_account_id"
        case destinationBeneficiaryID = "destination_beneficiary_id"
        case anchorTransferID = "anchor_transfer_id"
        case type
        case amount
        case fee
        case status
        case description
        case category
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}