package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kubuskotak/tyr"
	"github.com/rs/zerolog/log"
)

var sqlOpen = tyr.New

type SqlParams struct {
	Driver      string
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxLifeTime int
}

func PersistenceDB(args SqlParams, dbName string) (*tyr.Sql, func(), error) {
	conn, err := sqlOpen(tyr.SqlConnParams{
		Driver: args.Driver,
		Dsn:    args.DSN,
	})
	if err != nil {
		log.Error().Err(err).Msg("connection is failed")
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Db.Close()
	}

	conn.Db.SetConnMaxLifetime(time.Minute * time.Duration(args.MaxLifeTime))
	conn.Db.SetMaxOpenConns(args.MaxOpen)
	conn.Db.SetMaxIdleConns(args.MaxIdle)

	errChan := make(chan error)
	go Routine(context.Background(), dbName, conn.Db, errChan)

	return conn, cleanup, nil
}

func Routine(ctx context.Context, name string, db *sql.DB, er chan error) {
	retry := 0
	for retry < 2 {
		log.Info().Msgf("ping database... :%s", name)
		if db == nil {
			retry++
			time.Sleep(5 * time.Second)
			continue
		}
		if err := db.Ping(); err != nil {
			retry++
			log.Error().Err(err).Msgf("cannot connections to db %s, retry: %d", name, retry)
			time.Sleep(10 * time.Second)
			continue
		}
		er <- nil
		time.Sleep(100 * time.Second)
	}
	erPanic := fmt.Errorf("cannot re-connect db: %s, retry: %d", name, retry)
	log.Error().Err(erPanic).Msgf("cannot connection: %v to %s: %s, retry: %d", erPanic.Error(), "master::slave", name, retry)
	er <- erPanic
	close(er)
}
