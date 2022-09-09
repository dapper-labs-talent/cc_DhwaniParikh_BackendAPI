package database

import (
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"server/internal/config"
)

func Connect(config *config.Config, logger *zerolog.Logger) *pg.DB {
	var opt *pg.Options
	// Check to see if we have a database connection string
	if len(config.Database.ConnectionString) != 0 {
		o, err := pg.ParseURL(config.Database.ConnectionString)
		if err != nil {
			logger.Fatal().Err(err)
		}
		opt = o
	} else {
		logger.Fatal().Msg("No database configuration information supplied")
	}

	log.Info().Msg("Initializing database connection")
	db := pg.Connect(opt)

	log.Info().Msg("Established a successful connection!")

	return db
}
