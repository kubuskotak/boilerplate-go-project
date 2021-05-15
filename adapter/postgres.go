package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kubuskotak/boilerplate-go-project/pkg/db"
	"github.com/rs/zerolog/log"
)

var sqlOpen = db.New

func PostgresDB(args db.SqlParams) (*db.Sql, func(), error) {
	pg, err := sqlOpen(db.SqlConnParams{
		Driver: args.Driver,
		Dsn:    args.DSN,
	})
	if err != nil {

	}

	cleanup := func() {
		_ = pg.Db.Close()
	}

	pg.Db.SetMaxOpenConns(args.MaxOpen)
	pg.Db.SetMaxIdleConns(args.MaxIdle)
	pg.Db.SetConnMaxLifetime(time.Second * time.Duration(args.MaxLifeTime))

	errChan := make(chan error)
	go Routine(context.Background(), "Courier", pg.Db, errChan)

	return pg, cleanup, nil
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
