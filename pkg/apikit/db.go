package apikit

import (
	"context"
	"database/sql"
	"fmt"
)

func ConnectDBContext(ctx context.Context, cf DBConfig) (*sql.DB, error) {
	dsname := fmt.Sprintf("%s:%s@tcp(%s)/%s", cf.User, cf.Password, cf.Host, cf.Name)
	fmt.Printf("\nDatasource is %s\n", dsname)
	db, err := sql.Open("mysql", dsname)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
