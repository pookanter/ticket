package apikit

import (
	"context"
	"database/sql"
	"fmt"
	"ticket/pkg/db"
	"time"
)

type DBConfig struct {
	Host     string
	Name     string
	User     string
	Password string
	TimeOut  time.Duration
}

func ConnectDBContext(ctx context.Context, cf DBConfig) (*db.Queries, error) {
	dsname := fmt.Sprintf("%s:%s@%s?parseTime=true", cf.User, cf.Password, cf.Name)
	sqldb, err := sql.Open(cf.Host, dsname)
	if err != nil {
		return nil, err
	}

	err = sqldb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db.New(sqldb), nil
}
