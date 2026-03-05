-- =============================================================================
-- Migration: 000004_rename_code_add_encoded (DOWN)
-- Description: Reverts column rename and drops encoded_bookmark_code column.
-- =============================================================================

-- 1. Drop the encoded_bookmark_code column
ALTER TABLE bookmarks DROP COLUMN IF EXISTS encoded_bookmark_code;

-- 2. Rename 'sequence_code_id' back to 'code'
--    NOTE: PostgreSQL stores constraints using column OID (internal ID), not
--    column name strings. This means when a column is renamed, any constraint
--    referencing that column automatically follows the new name. So even if we
--    renamed the column AFTER creating the constraint, PostgreSQL would still
--    resolve it correctly. However, we rename FIRST for clarity, so step 3
--    can reference the final column name 'code' explicitly.
ALTER TABLE bookmarks RENAME COLUMN sequence_code_id TO code;

-- 3. Rename the unique constraint back (using final column name 'code')
ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS uni_sequence_code_id;
ALTER TABLE bookmarks ADD CONSTRAINT uni_code UNIQUE (code);
