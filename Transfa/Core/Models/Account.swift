/**
 * @description
 * This file defines the data model for a user's account (wallet).
 *
 * The `Account` struct represents a user's wallet within the Transfa system, which is backed by a
 * `DepositAccount` from the Anchor BaaS. It corresponds to a record in the `accounts` table.
 *
 * Key features:
 * - Conforms to `Codable` and `Identifiable` for API and SwiftUI integration.
 * - Uses an `AccountPurpose` enum to differentiate between wallet types (e.g., main vs. Money Drop).
 * - Implements `CodingKeys` to map between Swift's camelCase and the backend's snake_case.
 *
 * @dependencies
 * - Foundation: Provides `UUID` and `Date` types.
 */
import Foundation

/// Represents a user's wallet or special-purpose account within the Transfa application.
struct Account: Codable, Identifiable {
    /// The unique internal identifier for the account.
    let id: UUID
    
    /// The ID of the user who owns this account.
    let userID: UUID
    
    /// The unique identifier for the account from the BaaS provider (Anchor).
    let anchorAccountID: String
    
    /// The designated purpose of the account.
    let accountPurpose: AccountPurpose
    
    /// The current balance of the account, in the smallest currency unit (kobo).
    let balance: Int64
    
    /// The current status of the account (e.g., "active", "frozen").
    let status: String
    
    /// The timestamp when the account was created.
    let createdAt: Date
    
    /// The timestamp when the account was last updated.
    let updatedAt: Date
    
    /// Defines the possible purposes of an account. Raw values match the database enum.
    enum AccountPurpose: String, Codable {
        case mainWallet = "main_wallet"
        case moneyDropWallet = "money_drop_wallet"
    }
    
    /// Maps Swift's camelCase properties to the database's snake_case JSON keys.
    enum CodingKeys: String, CodingKey {
        case id
        case userID = "user_id"
        case anchorAccountID = "anchor_account_id"
        case accountPurpose = "account_purpose"
        case balance
        case status
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}