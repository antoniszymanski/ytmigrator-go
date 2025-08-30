// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/freetube"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type Cmd_yt2ft struct {
	YoutubeOptions
	FreetubeOptions
	common.ExportOptions
}

func (c *Cmd_yt2ft) Run() int {
	yt_client, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		common.Logger.Err(err).Msg("failed to create YouTube client")
		return 1
	}

	src := youtube.NewMigrator(yt_client)
	data, err := src.ExportFrom(c.ExportOptions)
	if err != nil {
		common.Logger.Err(err).Msg("failed to export data from YouTube")
		return 1
	}

	dst := freetube.NewMigrator(c.Dir)
	if err = dst.ImportTo(data); err != nil {
		common.Logger.Err(err).Msg("failed to import data to FreeTube")
		return 1
	}

	return 0
}
