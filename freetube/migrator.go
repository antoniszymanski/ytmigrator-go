// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package freetube

import (
	"github.com/antoniszymanski/ytmigrator-go/common"
)

type Migrator struct {
	dir string
}

var _ common.Migrator = Migrator{}

func NewMigrator(dir string) Migrator {
	return Migrator{dir: dir}
}
