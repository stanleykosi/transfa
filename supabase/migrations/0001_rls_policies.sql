/**
 * @description
 * Transfa App - Row Level Security (RLS) Policies
 *
 * This migration script enables Row Level Security on all tables containing
 * user-specific data and defines the policies that govern data access.
 * These policies are crucial for ensuring data privacy and security, preventing
 * users from accessing or modifying data that does not belong to them.
 *
 * The policies leverage Supabase's built-in `auth.uid()` function to identify
 * the currently authenticated user and check for ownership against foreign keys
 * in each table.
 *
 * RLS is enabled for the following tables:
 * - users
 * - user_settings
 * - accounts
 * - beneficiaries
 * - transactions
 * - subscriptions
 * - money_drops
 * - money_drop_claims
 * - payment_requests
 */

--==============================================================
-- RLS for `users` table
-- Users should be able to see and update their own profile.
--==============================================================
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can view their own profile."
ON public.users FOR SELECT
USING (auth.uid() = id);

CREATE POLICY "Users can update their own profile."
ON public.users FOR UPDATE
USING (auth.uid() = id);

--==============================================================
-- RLS for `user_settings` table
-- Users can perform all operations on their own settings.
--==============================================================
ALTER TABLE public.user_settings ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage their own settings."
ON public.user_settings FOR ALL
USING (auth.uid() = user_id);

--==============================================================
-- RLS for `accounts` table
-- Users can perform all operations on their own accounts.
--==============================================================
ALTER TABLE public.accounts ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can access their own accounts."
ON public.accounts FOR ALL
USING (auth.uid() = user_id);

--==============================================================
-- RLS for `beneficiaries` table
-- Users can perform all operations on their own beneficiaries.
--==============================================================
ALTER TABLE public.beneficiaries ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage their own beneficiaries."
ON public.beneficiaries FOR ALL
USING (auth.uid() = user_id);

--==============================================================
-- RLS for `transactions` table
-- Users can only view transactions where they are the sender or recipient.
-- Direct modification from the client is disabled.
--==============================================================
ALTER TABLE public.transactions ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can only see their own transactions."
ON public.transactions FOR SELECT
USING (auth.uid() = sender_user_id OR auth.uid() = recipient_user_id);

--==============================================================
-- RLS for `subscriptions` table
-- Users can perform all operations on their own subscription record.
--==============================================================
ALTER TABLE public.subscriptions ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage their own subscription."
ON public.subscriptions FOR ALL
USING (auth.uid() = user_id);

--==============================================================
-- RLS for `money_drops` table
-- Creators can perform all operations on their own money drops.
--==============================================================
ALTER TABLE public.money_drops ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Creators can manage their own money drops."
ON public.money_drops FOR ALL
USING (auth.uid() = creator_user_id);

--==============================================================
-- RLS for `money_drop_claims` table
-- Claimants can see their claims. Creators can see all claims on their drops.
--==============================================================
ALTER TABLE public.money_drop_claims ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Claimants can see their own claims."
ON public.money_drop_claims FOR SELECT
USING (auth.uid() = claimant_user_id);

CREATE POLICY "Creators can see claims on their money drops."
ON public.money_drop_claims FOR SELECT
USING (
  EXISTS (
    SELECT 1
    FROM public.money_drops md
    WHERE md.id = money_drop_claims.money_drop_id AND md.creator_user_id = auth.uid()
  )
);

--==============================================================
-- RLS for `payment_requests` table
-- Anyone can view payment requests. Only creators can manage them.
--==============================================================
ALTER TABLE public.payment_requests ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Anyone can view payment requests."
ON public.payment_requests FOR SELECT
USING (true);

CREATE POLICY "Creators can insert their own payment requests."
ON public.payment_requests FOR INSERT
WITH CHECK (auth.uid() = creator_user_id);

CREATE POLICY "Creators can update their own payment requests."
ON public.payment_requests FOR UPDATE
USING (auth.uid() = creator_user_id);

CREATE POLICY "Creators can delete their own payment requests."
ON public.payment_requests FOR DELETE
USING (auth.uid() = creator_user_id);