// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package freetube

type Migrator struct {
	dir string
}

func NewMigrator(dir string) *Migrator {
	return &Migrator{dir: dir}
}
