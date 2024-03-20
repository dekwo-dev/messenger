package database

import (
	"context"
	"strings"

	. "dekwo.dev/messager/logger"
	. "dekwo.dev/messager/env"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool = pool()

func pool() func() *pgxpool.Pool {
    const f = "pool"
    const file = "database/setup.go"

    var pool *pgxpool.Pool

    return func() *pgxpool.Pool {
        if pool != nil {
            return pool
        }

        url := GetEnv("POSTGRES_URL")
        if strings.Compare(url, "") == 0 {
            Fatal(50, file, f, "Postgres connection url is required in .env", nil) 
        }

        cfg, err := pgxpool.ParseConfig(url) 
        if err != nil {
            Fatal(50, file, f, "Failed to parse Postgres connection url", err)
        }

        if pool, err = pgxpool.NewWithConfig(context.Background(), cfg); err != nil {
            Fatal(50, file, f, "Failed to create a new pgxpool.Pool", err)
        }

        return pool
    }
}
