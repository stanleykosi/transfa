/**
 * @description
 * Transfa App - Initial Database Schema
 * 
 * This script sets up the complete initial database schema for the Transfa application.
 * It includes all tables, relationships, constraints, and helper functions required
 * for the MVP as defined in the technical specification.
 * 
 * Key Features:
 * - Creates all core tables for users, accounts, transactions, and features.
 * - Establishes foreign key relationships to ensure data integrity.
 * - Implements CHECK constraints for enumerated types (e.g., statuses, types).
 * - Uses UUIDs for all primary keys for scalability.
 * - Includes a trigger function to automatically update `updated_at` timestamps on all tables.
 *
 * Migration Order (to respect dependencies):
 * 1. trigger_set_timestamp() helper function
 * 2. users
 * 3. accounts
 * 4. beneficiaries
 * 5. user_settings
 * 6. subscriptions
 * 7. transactions
 * 8. money_drops
 * 9. money_drop_claims
 * 10. payment_requests
 */

--==============================================================
-- HELPER FUNCTION & TRIGGERS
--==============================================================

--
-- Name: trigger_set_timestamp(); Type: FUNCTION;
-- Description: A trigger function to automatically set the updated_at column to the current timestamp.
--
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


--==============================================================
-- TABLE CREATION
--==============================================================

--
-- Table: users
-- Description: Stores public user profile information, linking to Supabase auth and Anchor BaaS.
--
CREATE TABLE public.users (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_id text NOT NULL UNIQUE,
    username text NOT NULL UNIQUE,
    account_type text NOT NULL CHECK (account_type IN ('personal', 'merchant')),
    anchor_customer_id text NOT NULL UNIQUE,
    kyc_status text NOT NULL DEFAULT 'pending' CHECK (kyc_status IN ('pending', 'approved', 'rejected')),
    profile_image_url text,
    allow_sending boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_id_fkey FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);
COMMENT ON TABLE public.users IS 'Stores user profile information and links to authentication and BaaS providers.';

-- Add trigger for users table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: accounts
-- Description: Stores user bank accounts provided by the BaaS (Anchor).
--
CREATE TABLE public.accounts (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    anchor_account_id text NOT NULL UNIQUE,
    account_purpose text NOT NULL CHECK (account_purpose IN ('main_wallet', 'money_drop_wallet')),
    balance bigint NOT NULL DEFAULT 0,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.accounts IS 'Represents user wallets and other special-purpose accounts within Transfa.';


-- Add trigger for accounts table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.accounts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: beneficiaries
-- Description: Stores external bank accounts that users have added for transfers.
--
CREATE TABLE public.beneficiaries (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    anchor_counterparty_id text NOT NULL UNIQUE,
    account_name text NOT NULL,
    account_number text NOT NULL,
    bank_name text NOT NULL,
    bank_code text NOT NULL,
    is_default boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.beneficiaries IS 'External bank accounts saved by users for withdrawals (Self Transfers).';


-- Add trigger for beneficiaries table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.beneficiaries
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: user_settings
-- Description: Stores user-specific settings and preferences.
--
CREATE TABLE public.user_settings (
    user_id uuid NOT NULL PRIMARY KEY REFERENCES public.users(id) ON DELETE CASCADE,
    default_beneficiary_id uuid REFERENCES public.beneficiaries(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.user_settings IS 'Stores user-specific preferences, such as default accounts.';


-- Add trigger for user_settings table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.user_settings
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: subscriptions
-- Description: Manages user subscription status and entitlements.
--
CREATE TABLE public.subscriptions (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL UNIQUE REFERENCES public.users(id) ON DELETE CASCADE,
    status text NOT NULL DEFAULT 'free' CHECK (status IN ('free', 'active', 'past_due', 'cancelled')),
    auto_renew boolean NOT NULL DEFAULT false,
    current_period_ends_at timestamptz,
    monthly_external_transfers_used integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.subscriptions IS 'Tracks user subscription tiers, payment status, and usage.';


-- Add trigger for subscriptions table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.subscriptions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: transactions
-- Description: A central log of all money movements within the application.
--
CREATE TABLE public.transactions (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_user_id uuid REFERENCES public.users(id),
    recipient_user_id uuid REFERENCES public.users(id),
    source_account_id uuid REFERENCES public.accounts(id),
    destination_account_id uuid REFERENCES public.accounts(id),
    destination_beneficiary_id uuid REFERENCES public.beneficiaries(id),
    anchor_transfer_id text UNIQUE,
    type text NOT NULL CHECK (type IN ('p2p', 'self_transfer', 'money_drop_funding', 'money_drop_claim', 'subscription_fee', 'wallet_funding')),
    amount bigint NOT NULL,
    fee bigint NOT NULL DEFAULT 0,
    status text NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'reversed')),
    description text,
    category text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.transactions IS 'Records all financial transactions occurring on the platform.';


-- Add trigger for transactions table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.transactions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: money_drops
-- Description: Stores the configuration and state of a Money Drop instance.
--
CREATE TABLE public.money_drops (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_user_id uuid NOT NULL REFERENCES public.users(id),
    funding_account_id uuid NOT NULL REFERENCES public.accounts(id),
    total_amount bigint NOT NULL,
    amount_per_claim bigint NOT NULL,
    total_claims_allowed integer NOT NULL,
    claims_made_count integer NOT NULL DEFAULT 0,
    status text NOT NULL CHECK (status IN ('active', 'completed', 'expired')),
    expiry_timestamp timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.money_drops IS 'Defines a Money Drop event, including its rules and funding.';


-- Add trigger for money_drops table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.money_drops
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


--
-- Table: money_drop_claims
-- Description: Records which users have claimed from a specific Money Drop.
--
CREATE TABLE public.money_drop_claims (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    money_drop_id uuid NOT NULL REFERENCES public.money_drops(id) ON DELETE CASCADE,
    claimant_user_id uuid NOT NULL REFERENCES public.users(id),
    claimed_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (money_drop_id, claimant_user_id)
);
COMMENT ON TABLE public.money_drop_claims IS 'Tracks individual claims against a Money Drop to prevent duplicates.';


--
-- Table: payment_requests
-- Description: Stores details of payment requests created by users.
--
CREATE TABLE public.payment_requests (
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_user_id uuid NOT NULL REFERENCES public.users(id),
    amount bigint NOT NULL,
    description text,
    image_url text,
    status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'fulfilled')),
    fulfilled_at timestamptz,
    fulfilled_by_transaction_id uuid REFERENCES public.transactions(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE public.payment_requests IS 'Stores details for user-generated payment requests.';


-- Add trigger for payment_requests table
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON public.payment_requests
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();