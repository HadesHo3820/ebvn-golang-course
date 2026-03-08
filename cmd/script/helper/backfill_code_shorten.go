package helper

import (
	"errors"
	"fmt"

	"github.com/HadesHo3820/ebvn-golang-course/pkg/base62"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// For Pluralized Table Name, please reference to https://gorm.io/docs/conventions.html#Pluralized-Table-Name

type bookmark struct {
	ID             string `gorm:"column:id"`
	SequenceCodeID int64  `gorm:"column:sequence_code_id"`
}

// You can change the default table name by implementing the Tabler interface
func (b bookmark) TableName() string {
	return "bookmarks"
}

type schemaMigration struct {
	Version int64 `gorm:"column:version"`
}

func (s schemaMigration) TableName() string {
	return "schema_migrations"
}

// BackfillForEncodedBookmarkCodeCol backfills the encoded_bookmark_code column for existing bookmarks.
//
// This function retrieves all bookmarks that have a NULL encoded_bookmark_code (ordered by creation time),
// encodes each bookmark's sequence_code_id to Base62, and updates the encoded_bookmark_code column.
//
// Prerequisites:
//   - Migration 000004 (which adds the encoded_bookmark_code column) must have been applied.
//
// Idempotency:
//   - Safe to re-run. Only processes rows where encoded_bookmark_code IS NULL.
func BackfillForEncodedBookmarkCodeCol(db *gorm.DB) error {
	// 1. Guard: ensure migration 000004 has been applied
	sm := &schemaMigration{}
	if err := db.Model(&schemaMigration{}).Where("version = ?", 4).First(sm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Migration 000004 has not been applied yet. Skipping backfill.")
			return fmt.Errorf("migration 000004 not applied: %w", err)
		}
		log.Error().Err(err).Msg("Failed to query schema_migrations table.")
		return err
	}

	// 2. Fetch all bookmarks that need backfilling
	log.Info().Msg("Starting backfill for encoded_bookmark_code column in bookmarks table.")
	var bookmarks []bookmark
	if err := db.Model(&bookmark{}).Where("encoded_bookmark_code IS NULL").Order("created_at ASC").Find(&bookmarks).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch bookmarks for backfilling encoded_bookmark_code column.")
		return err
	}

	if len(bookmarks) == 0 {
		log.Info().Msg("No bookmarks found for backfilling encoded_bookmark_code column.")
		return nil
	}

	// 3. Backfill within a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, b := range bookmarks {
			// Encode the actual sequence_code_id value (not an arbitrary loop index)
			encoded, err := base62.Encode(b.SequenceCodeID)
			if err != nil {
				return fmt.Errorf("failed to encode sequence_code_id %d for bookmark %s: %w", b.SequenceCodeID, b.ID, err)
			}

			if err := tx.Model(&bookmark{}).Where("id = ?", b.ID).Update("encoded_bookmark_code", encoded).Error; err != nil {
				return fmt.Errorf("failed to update bookmark %s: %w", b.ID, err)
			}
			log.Info().Str("id", b.ID).Int64("sequence_code_id", b.SequenceCodeID).Str("encoded", encoded).Msg("Updated encoded_bookmark_code")
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Transaction failed while backfilling encoded_bookmark_code column.")
		return err
	}

	log.Info().Int("count", len(bookmarks)).Msg("Successfully backfilled encoded_bookmark_code column for all bookmarks.")
	return nil
}
