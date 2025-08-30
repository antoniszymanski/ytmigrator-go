// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tubular

import (
	"database/sql"

	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/tubular/internal"
	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	queries *internal.Queries
	db      *sql.DB
}

var _ common.Migrator = (*Migrator)(nil)

func (m *Migrator) Close() error {
	return m.db.Close()
}

func NewMigrator(dsn string) (*Migrator, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &Migrator{queries: internal.New(db), db: db}, nil
}
