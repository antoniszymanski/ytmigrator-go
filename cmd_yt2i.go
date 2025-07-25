// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"time"

	invidiousapi "github.com/antoniszymanski/invidious-go"
	"github.com/antoniszymanski/ytmigrator-go/invidious"
	"github.com/antoniszymanski/ytmigrator-go/youtube"
)

type cmd_yt2i struct {
	youtubeConfig   `prefix:"yt-"`
	invidiousConfig `prefix:"i-"`
}

func (c *cmd_yt2i) Run() error {
	ytClient, err := youtube.NewService(c.Credentials, c.Token)
	if err != nil {
		return err
	}
	src := youtube.NewMigrator(ytClient)

	iClient := invidiousapi.NewClient(c.InstanceURL)
	err = iClient.AuthorizeToken(invidiousapi.AuthorizeTokenRequest{
		Scopes: []string{":subscriptions*", ":playlist*"},
		Expire: time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		return err
	}

	takeoutFile, err := os.OpenFile(c.Takeout, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	dst, err := invidious.NewMigrator(takeoutFile, iClient)
	if err != nil {
		return err
	}
	defer dst.Close() //nolint:errcheck

	data, err := src.Export(cli.ExportOptions)
	if err != nil {
		return err
	}
	return dst.Import(data)
}
