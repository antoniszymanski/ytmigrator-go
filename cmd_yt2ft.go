// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/antoniszymanski/ytmigrator-go/freetube"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type Cmd_yt2ft struct {
	YoutubeOptions
	FreetubeOptions
}

func (c *Cmd_yt2ft) Run() int {
	ytClient, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		logger.Err(err).Msg("failed to create YouTube service")
		return 1
	}
	src := youtube.NewMigrator(ytClient)
	dst := freetube.NewMigrator(c.Dir)

	data, err := src.Export(args.ExportOptions)
	if err != nil {
		logger.Err(err).Msg("failed to export data from YouTube")
		return 1
	}

	if err = dst.Import(data); err != nil {
		logger.Err(err).Msg("failed to import data to FreeTube")
		return 1
	}

	return 0
}
