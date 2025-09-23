package main

import (
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/tubular"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type Cmd_yt2t struct {
	YoutubeOptions
	TubularOptions
	common.ExportOptions
}

func (c *Cmd_yt2t) Run() int {
	yt_client, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		common.Logger.Err(err).Msg("failed to create YouTube client")
		return 1
	}

	dst, err := tubular.NewMigrator(c.DSN)
	if err != nil {
		common.Logger.Err(err).Msg("failed to create Tubular migrator")
		return 1
	}
	defer dst.Close() //nolint:errcheck

	src := youtube.NewMigrator(yt_client)
	data, err := src.ExportFrom(c.ExportOptions)
	if err != nil {
		common.Logger.Err(err).Msg("failed to export data from YouTube")
		return 1
	}

	if err = dst.ImportTo(data); err != nil {
		common.Logger.Err(err).Msg("failed to import data to Tubular")
		return 1
	}

	return 0
}
