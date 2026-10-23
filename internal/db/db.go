package db

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func Open(ctx context.Context, path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?"+
		"_pragma=journal_mode(WAL)&"+
		"_pragma=busy_timeout(5000)&"+
		"_pragma=foreign_keys(1)&"+
		"_pragma=synchronous(NORMAL)&"+
		"_txlock=immediate", path)
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(4)
	db.SetConnMaxIdleTime(0)
	err = db.PingContext(ctx)
	if err != nil {
		zap.L().Error("ping to database failed")
		return nil, err
	}
	return db, nil
}
