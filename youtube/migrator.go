// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package youtube

import (
	"github.com/antoniszymanski/ytmigrator-go/common"
	"google.golang.org/api/youtube/v3"
)

type Migrator struct {
	client *youtube.Service
}

var _ common.Migrator = (*Migrator)(nil)

func NewMigrator(client *youtube.Service) *Migrator {
	return &Migrator{client: client}
}
