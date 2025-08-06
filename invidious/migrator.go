// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"io"
	"os"

	"github.com/antoniszymanski/invidious-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/invidious/models"
	"github.com/dsnet/try"
	"github.com/go-json-experiment/json"
	"github.com/rs/zerolog"
)

type Migrator struct {
	logger      *zerolog.Logger
	takeoutFile *os.File
	takeout     models.Takeout
	client      *invidious.Client
}

var _ common.Migrator = (*Migrator)(nil)

func (m *Migrator) SetLogger(logger *zerolog.Logger) {
	m.logger = logger
}

func (m *Migrator) Close() (err error) {
	defer try.Handle(&err)

	try.E1(m.takeoutFile.Seek(0, io.SeekStart))
	try.E(m.takeoutFile.Truncate(0))
	try.E(json.MarshalWrite(m.takeoutFile, &m.takeout))
	return m.takeoutFile.Close()
}

func NewMigrator(takeoutFile *os.File, client *invidious.Client) (*Migrator, error) {
	var takeout models.Takeout
	if err := json.UnmarshalRead(takeoutFile, &takeout); err != nil {
		return nil, err
	}
	return &Migrator{
		takeoutFile: takeoutFile,
		takeout:     takeout,
		client:      client,
	}, nil
}
