package pg

import "github.com/jackc/pgx/v4/pgxpool"

func NewDB(s string) (*pgxpool.Pool, error)
