package main

import (
	"coding-challenge-go/pkg/api"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	db, err := sql.Open("mysql", "user:password@tcp(db:3306)/product")

	if err != nil {
		log.Error().Err(err).Msg("Fail to create server")
		return
	}

	defer db.Close()

	engine, err := api.CreateAPIEngine(db)

	if err != nil {
		log.Error().Err(err).Msg("Fail to create server")
		return
	}

	log.Info().Msg("Start server")
	log.Fatal().Err(engine.Run(os.Getenv("LISTEN"))).Msg("Fail to listen and serve")
}
