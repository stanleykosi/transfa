-- NOTE: This must be executed by the table owner (supabase_storage_admin) or a superuser/service role connection.

-- Public read for bucket
CREATE POLICY "public_read_user_content"
ON storage.objects FOR SELECT
USING (bucket_id = 'user_content');

-- Insert only into user's top-level folder (uses generated path_tokens)
CREATE POLICY "authenticated_insert_own_folder"
ON storage.objects FOR INSERT
TO authenticated
WITH CHECK (
  bucket_id = 'user_content' AND
  (path_tokens)[1] = (SELECT auth.uid())::text
);

-- Update only own files (owner is uuid column)
CREATE POLICY "authenticated_update_own_files"
ON storage.objects FOR UPDATE
TO authenticated
USING (
  bucket_id = 'user_content' AND
  (SELECT auth.uid())::uuid = owner
)
WITH CHECK (
  bucket_id = 'user_content' AND
  (SELECT auth.uid())::uuid = owner
);

-- Delete only own files
CREATE POLICY "authenticated_delete_own_files"
ON storage.objects FOR DELETE
TO authenticated
USING (
  bucket_id = 'user_content' AND
  (SELECT auth.uid())::uuid = owner
);

-- Helpful SELECT to verify current role when running
SELECT current_user AS current_role, session_user AS session_user;