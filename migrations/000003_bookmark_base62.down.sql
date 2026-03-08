-- =============================================================================
-- Migration: 000003_bookmark_base62 (DOWN)
-- Description: Reverts bookmark code from bigint back to varchar
-- =============================================================================

-- 1. Drop the unique constraint
ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS uni_code;

-- 2. Remove the default sequence value
ALTER TABLE bookmarks ALTER COLUMN code DROP DEFAULT;

-- 3. Drop the sequence
DROP SEQUENCE IF EXISTS bookmarks_code_seq;

-- 4. Change the code column back to varchar
ALTER TABLE bookmarks ALTER COLUMN code TYPE varchar(10) USING code::varchar;

-- 5. Re-add the unique constraint
ALTER TABLE bookmarks ADD CONSTRAINT uni_code UNIQUE (code);
