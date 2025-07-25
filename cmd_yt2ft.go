// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/antoniszymanski/ytmigrator-go/freetube"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type cmd_yt2ft struct {
	youtubeConfig  `prefix:"yt-"`
	freetubeConfig `prefix:"ft-"`
}

func (c *cmd_yt2ft) Run() error {
	ytClient, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		return err
	}
	src := youtube.NewMigrator(ytClient)
	dst := freetube.NewMigrator(c.Dir)

	data, err := src.Export(cli.ExportOptions)
	if err != nil {
		return err
	}
	return dst.Import(data)
}
