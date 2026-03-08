package main

import (
	"context"

	"github.com/HadesHo3820/ebvn-golang-course/cmd/script/helper"
	"github.com/HadesHo3820/ebvn-golang-course/internal/infrastructure"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	dbClient := infrastructure.CreateSQLDB()
	err := helper.BackfillForEncodedBookmarkCodeCol(dbClient)
	if err != nil {
		// Fatal logs the error and calls os.Exit(1), ensuring Docker/K8s
		// detects the failure via a non-zero exit code.
		log.Fatal().Err(err).Msg("Failed to backfill encoded_bookmark_code column")
	}

	log.Info().Msg("Clearing Redis cache after backfill.")
	redisClient := infrastructure.CreateRedisCacheConn()
	err = redisClient.FlushAll(ctx).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to clear Redis cache after backfill")
	}

	log.Info().Msg("Successfully backfilled encoded_bookmark_code column and cleared Redis cache.")
}
