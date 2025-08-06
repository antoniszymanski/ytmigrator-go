// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package freetube

import (
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/rs/zerolog"
)

type Migrator struct {
	logger *zerolog.Logger
	dir    string
}

var _ common.Migrator = (*Migrator)(nil)

func (m *Migrator) SetLogger(logger *zerolog.Logger) {
	m.logger = logger
}

func NewMigrator(dir string) *Migrator {
	return &Migrator{dir: dir}
}
