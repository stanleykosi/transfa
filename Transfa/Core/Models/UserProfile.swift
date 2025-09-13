/**
 * @description
 * This file defines the core data model for a user profile in the Transfa app.
 *
 * The `UserProfile` struct is the Swift representation of a user's data as stored in the
 * `users` table of the Supabase database. It includes essential information like the user's
 * unique identifiers, username, account type, and settings.
 *
 * Key features:
 * - Conforms to `Codable` for seamless JSON serialization/deserialization with the backend.
 * - Conforms to `Identifiable`, making it compatible with SwiftUI lists and views.
 * - Uses an `AccountType` enum for type-safe handling of user roles.
 * - Implements `CodingKeys` to map between Swift's camelCase properties and the database's snake_case columns.
 *
 * @dependencies
 * - Foundation: Provides `UUID` and `Date` types.
 */
import Foundation

/// Represents a user's profile within the Transfa application.
struct UserProfile: Codable, Identifiable {
    /// The unique identifier for the user, matching the Supabase `auth.users` ID.
    let id: UUID
    
    /// The user's unique, publicly searchable username.
    let username: String
    
    /// The type of account the user has, either personal or merchant.
    let accountType: AccountType
    
    /// The URL for the user's profile picture, stored in Supabase Storage. Optional.
    let profileImageURL: String?
    
    /// A flag indicating if the user's account is allowed to send money. Defaults to `false` for merchants.
    let allowSending: Bool
    
    /// The current status of the user's KYC (Know Your Customer) verification.
    let kycStatus: String
    
    /// The timestamp when the user profile was created.
    let createdAt: Date
    
    /// The timestamp when the user profile was last updated.
    let updatedAt: Date
    
    /// Defines the possible types of user accounts. Raw values match the database enum.
    enum AccountType: String, Codable {
        case personal
        case merchant
    }
    
    /// Maps Swift's camelCase properties to the database's snake_case JSON keys.
    enum CodingKeys: String, CodingKey {
        case id
        case username
        case accountType = "account_type"
        case profileImageURL = "profile_image_url"
        case allowSending = "allow_sending"
        case kycStatus = "kyc_status"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}