-- =============================================================================
-- Migration: 000004_rename_code_add_encoded
-- Description: Renames 'code' column to 'sequence_code_id' for naming
--              consistency, and adds 'encoded_bookmark_code' column to store
--              the Base62-encoded value directly in the database.
-- =============================================================================

-- 1. Rename 'code' → 'sequence_code_id'
--    This ensures consistent naming across the schema. The column still holds
--    the auto-incremented bigint value from the bookmarks_code_seq sequence.
ALTER TABLE bookmarks RENAME COLUMN code TO sequence_code_id;

-- 2. Rename the unique constraint to match the new column name
ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS uni_code;
ALTER TABLE bookmarks ADD CONSTRAINT uni_sequence_code_id UNIQUE (sequence_code_id);

-- 3. Add 'encoded_bookmark_code' column (nullable)
--    Stores the Base62-encoded string of sequence_code_id.
--    Default is NULL; populated by the application during bookmark creation.
ALTER TABLE bookmarks ADD COLUMN encoded_bookmark_code VARCHAR(20) DEFAULT NULL;
