// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"time"

	invidiousapi "github.com/antoniszymanski/invidious-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/invidious"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type Cmd_yt2i struct {
	YoutubeOptions
	InvidiousOptions
	common.ExportOptions
}

func (c *Cmd_yt2i) Run() int {
	ytClient, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		logger.Err(err).Msg("failed to create YouTube service")
		return 1
	}
	src := youtube.NewMigrator(ytClient)

	iClient := invidiousapi.NewClient(c.InstanceURL)
	err = iClient.AuthorizeToken(invidiousapi.AuthorizeTokenRequest{
		Scopes: []string{":subscriptions*", ":playlist*"},
		Expire: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		logger.Err(err).Msg("failed to authorize Invidious token")
		return 1
	}

	takeoutFile, err := os.OpenFile(c.Takeout, os.O_RDWR, 0600)
	if err != nil {
		logger.Err(err).Msg("failed to open takeout file")
		return 1
	}
	dst, err := invidious.NewMigrator(takeoutFile, iClient)
	if err != nil {
		logger.Err(err).Msg("failed to create Invidious migrator")
		return 1
	}
	defer dst.Close() //nolint:errcheck

	data, err := src.Export(c.ExportOptions)
	if err != nil {
		logger.Err(err).Msg("failed to export data from YouTube")
		return 1
	}

	if err = dst.Import(data); err != nil {
		logger.Err(err).Msg("failed to import data to Invidious")
		return 1
	}

	return 0
}
