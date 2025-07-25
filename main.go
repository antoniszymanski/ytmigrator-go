// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/alecthomas/kong"
	"github.com/antoniszymanski/ytmigrator-go/common"
)

var cli struct {
	Cmd_yt2i  *cmd_yt2i  `cmd:"" name:"yt2i"`
	Cmd_yt2ft *cmd_yt2ft `cmd:"" name:"yt2ft"`
	common.ExportOptions
}

type youtubeConfig struct {
	Credentials string `type:"existingfile" default:"credentials.json"`
	Token       string `type:"path" default:"token.json"`
}

type invidiousConfig struct {
	Takeout     string `type:"existingfile" default:"subscription_manager.json"`
	InstanceURL string `required:""`
}

type freetubeConfig struct {
	Dir string `type:"path" default:"freetube"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("ytmigrator"),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
