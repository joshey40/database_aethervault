package services

import (
	"database/sql"

	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
)

type SyncService struct {
	DB      *sql.DB
	Queries *dbqueries.Queries
}

func NewSyncService(db *sql.DB, queries *dbqueries.Queries) *SyncService {
	return &SyncService{
		DB:      db,
		Queries: queries,
	}
}
