package connection

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Zhiyenbek/sp-positions-main-service/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgresDB(cfg *config.DBConf) (*pgxpool.Pool, error) {
	dbURI, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Println("Couldn't get database url. Continuing with config")
		dbURI = fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	} else {
		log.Println("Database url: ", dbURI)
	}
	log.Println(dbURI)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.TimeOut)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	c, ioErr := ioutil.ReadFile("./scripts/init.sql")
	if ioErr != nil {
		log.Println(ioErr)
	} else {
		sql := string(c)
		_, err = pool.Exec(ctx, sql)
		if err != nil {
			log.Println(err)
		}
	}
	return pool, nil
}
