package model

// Bookmark represents a shortened URL bookmark in the system.
// This struct maps to the "bookmarks" table in the database and stores
// URL shortening information with ownership tracking and soft delete support.
//
// Each bookmark is owned by a single user and will be automatically deleted
// when its owner is deleted (ON DELETE CASCADE constraint in the database).
//
// Fields:
//   - Base: Embedded struct providing ID, CreatedAt, UpdatedAt, and DeletedAt
//   - Description: Optional user-provided description or title for the bookmark
//   - URL: The original long URL that the short code redirects to
//   - Code: The unique short code used for redirection (e.g., "abc123")
//   - UserID: Foreign key referencing the user who created this bookmark
//   - User: The associated User object (excluded from JSON, loaded via GORM association)
type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code" gorm:"unique"`
	UserID      string `json:"user_id"`
	User        *User  `gorm:"references:ID" json:"-"`
}

// User represents the "Belongs To" relationship with the User model.
//
// The Mechanism:
// In this struct:
//   UserID string // <--- SOURCE (Foreign Key)
//   User   User   `gorm:"references:ID"` // <--- TARGET (Relationship)
//
// 1. Defaults: Because the relationship field is named User, GORM automatically looks for
//    a field named User + ID = UserID in this same struct to use as the value source.
//
// 2. Execution:
//    - GORM reads the value from Bookmark.UserID (e.g., "123").
//    - It runs a query on the users table: SELECT * FROM users WHERE id = "123"
//      (because references:ID points to the ID column in User).
//
// If you named it differently...
// If your source field was named something else, like OwnerID, GORM wouldn't find it
// automatically. You would need to tell it explicitly using foreignKey:
//
//   type Bookmark struct {
//       OwnerID string // Custom name
//       User    User   `gorm:"foreignKey:OwnerID;references:ID"`
//                                         ^ Use 'OwnerID' as the source value
//   }
//
// Summary:
//   - foreignKey: "Which field in THIS struct holds the value?" (Default: UserID)
//   - references: "Which field in the OTHER struct should we match against?" (Default: ID)
