-- =============================================================================
-- Migration: 000002_add_bookmark
-- Description: Creates the bookmarks table for storing shortened URL bookmarks
-- =============================================================================
-- This table stores URL bookmarks with unique short codes for redirection.
-- Each bookmark is owned by a user and supports soft deletion via deleted_at.
-- =============================================================================

CREATE TABLE bookmarks 
(
    -- Primary key: UUID stored as string
    id varchar(36) unique,
    
    -- Optional description/title for the bookmark
    description varchar(256),
    
    -- Original URL to redirect to (required)
    url varchar(2048) not null,
    
    -- Short code for the bookmark (e.g., "abc123" for /links/redirect/abc123)
    code varchar(10) not null,
    
    -- Foreign key: References the user who created this bookmark
    user_id varchar(36) not null,
    
    -- Timestamps for auditing (created_at and updated_at auto-managed)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Soft deletion timestamp (NULL means not deleted)
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints:
    CONSTRAINT bookmark_pkey PRIMARY KEY (id),           -- Ensures id is unique and indexed
    CONSTRAINT uni_code UNIQUE (code),                   -- Ensures short codes are unique across all bookmarks
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)          -- Links bookmark to user with CASCADE deletion
        REFERENCES users (id) ON DELETE CASCADE          -- ON DELETE CASCADE: When a user is deleted,
                                                         -- all their bookmarks are automatically deleted too.
                                                         -- This prevents orphaned bookmarks and maintains data consistency.
);
	