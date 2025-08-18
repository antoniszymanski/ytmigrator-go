// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tubular

import (
	"database/sql"

	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/tubular/internal"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

type Migrator struct {
	logger  *zerolog.Logger
	queries *internal.Queries
	db      *sql.DB
}

var _ common.Migrator = (*Migrator)(nil)

func (m *Migrator) SetLogger(logger *zerolog.Logger) {
	m.logger = logger
}

func (m *Migrator) Close() error {
	return m.db.Close()
}

func New(dsn string) (*Migrator, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &Migrator{queries: internal.New(db), db: db}, nil
}
