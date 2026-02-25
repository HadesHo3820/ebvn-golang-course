-- =============================================================================
-- Migration: 000003_bookmark_base62
-- Description: Changes bookmark code from varchar to bigint with auto-increment
--              to support base62 encoding for short URL codes.
-- =============================================================================

-- 1. Drop the unique constraint on the code column
--    Before we can change the column type, we must remove any indexes or 
--    constraints that depend on the existing varchar column.
ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS uni_code;

-- 2. Create the sequence FIRST (we need it in step 3)
--    A sequence is an independent database object that generates unique, 
--    strictly increasing numbers (1, 2, 3...). This is what provides our 
--    collision-free IDs for the base62 encoding.
CREATE SEQUENCE IF NOT EXISTS bookmarks_code_seq;

-- 3. Change the code column type from varchar to bigint
--    We use nextval('bookmarks_code_seq') so each existing row gets a unique
--    sequential code (1, 2, 3...) instead of all being set to the same value.
--    This prevents unique constraint violations when re-adding the constraint.
ALTER TABLE bookmarks ALTER COLUMN code TYPE bigint USING nextval('bookmarks_code_seq');

-- 4. Set the default value for the code column to use the sequence
--    This binds the sequence to the column. Now, whenever a new bookmark 
--    is INSERTed without explicitly providing a 'code', PostgreSQL will 
--    automatically fetch the next number from the sequence and use it.
ALTER TABLE bookmarks ALTER COLUMN code SET DEFAULT nextval('bookmarks_code_seq');

-- 5. Re-add the unique constraint on the code column
--    We re-apply the unique constraint to the new bigint column to ensure 
--    data integrity and to recreate the unique index, which makes querying 
--    by code extremely fast (O(1) lookups for the redirect endpoint).
ALTER TABLE bookmarks ADD CONSTRAINT uni_code UNIQUE (code);

