/**
 * @description
 * This file defines the data model for a user's beneficiary.
 *
 * A `Beneficiary` represents an external bank account that a user has saved for the purpose
 * of "Self Transfers" (withdrawals). This model corresponds to a record in the `beneficiaries`
 * table and is linked to a `CounterParty` resource in the Anchor BaaS.
 *
 * Key features:
 * - Conforms to `Codable` for API integration and `Identifiable` for use in SwiftUI lists.
 * - Contains all necessary details to identify an external bank account.
 * - Implements `CodingKeys` to map between Swift's camelCase and the backend's snake_case.
 *
 * @dependencies
 * - Foundation: Provides `UUID` and `Date` types.
 */
import Foundation

/// Represents a user's saved external bank account for withdrawals.
struct Beneficiary: Codable, Identifiable {
    /// The unique internal identifier for the beneficiary record.
    let id: UUID
    
    /// The ID of the user who owns this beneficiary record.
    let userID: UUID
    
    /// The unique identifier for the counterparty from the BaaS provider (Anchor).
    let anchorCounterpartyID: String
    
    /// The name of the account holder for the external bank account.
    let accountName: String
    
    /// The account number of the external bank account.
    let accountNumber: String
    
    /// The name of the bank (e.g., "Guaranty Trust Bank").
    let bankName: String
    
    /// The unique code for the bank.
    let bankCode: String
    
    /// A flag indicating if this is the user's default beneficiary for receiving payments.
    let isDefault: Bool
    
    /// The timestamp when the beneficiary was created.
    let createdAt: Date
    
    /// The timestamp when the beneficiary was last updated.
    let updatedAt: Date
    
    /// Maps Swift's camelCase properties to the database's snake_case JSON keys.
    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case anchorCounterpartyID = "anchor_counterparty_id"
        case accountName = "account_name"
        case accountNumber = "account_number"
        case bankName = "bank_name"
        case bankCode = "bank_code"
        case isDefault = "is_default"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}