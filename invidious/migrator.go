// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"io"
	"os"

	"github.com/antoniszymanski/invidious-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/invidious/models"
	"github.com/go-json-experiment/json"
)

type Migrator struct {
	takeoutFile *os.File
	takeout     models.Takeout
	client      *invidious.Client
}

var _ common.Migrator = (*Migrator)(nil)

func (m *Migrator) Close() error {
	if _, err := m.takeoutFile.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := m.takeoutFile.Truncate(0); err != nil {
		return err
	}
	if err := json.MarshalWrite(m.takeoutFile, &m.takeout); err != nil {
		return err
	}
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
