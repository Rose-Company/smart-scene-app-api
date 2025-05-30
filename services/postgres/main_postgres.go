package postgres

import (
	"log"
	"smart-scene-app-api/config"
	postgres2 "smart-scene-app-api/pkg/postgres"

	"gorm.io/gorm/logger"
)

func NewMainPostgres(prefix string) (error, *postgres2.Postgres) {
	var debugMode logger.LogLevel
	mode := config.Config.Postgres.GormDebug
	if mode == "release" {
		debugMode = logger.Silent
	} else {
		debugMode = logger.Info
	}

	postgresParams := postgres2.ConfigureParams{
		User:      config.Config.Postgres.User,
		Password:  config.Config.Postgres.Pass,
		Host:      config.Config.Postgres.Host,
		Port:      config.Config.Postgres.Port,
		Params:    config.Config.Postgres.Params,
		Database:  config.Config.Postgres.Db,
		DebugMode: debugMode,
	}

	log.Printf("Connecting to Postgres URI: %s\n", postgres2.GetPostgresUri(postgresParams))

	pg := postgres2.Postgres{}
	err := pg.Configure(prefix, postgresParams)
	if err != nil {
		return err, nil
	}

	if err := pg.Run(); err != nil {
		return err, nil
	}

	return nil, &pg
}
