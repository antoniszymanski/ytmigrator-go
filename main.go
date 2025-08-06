// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"

	"github.com/alexflint/go-arg"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/rs/zerolog"
)

var args struct {
	*Cmd_yt2i  `arg:"subcommand:yt2i"`
	*Cmd_yt2ft `arg:"subcommand:yt2ft"`
	common.ExportOptions
}

type YoutubeOptions struct {
	Credentials string `arg:"--youtube.credentials" default:"credentials.json"`
	Token       string `arg:"--youtube.token" default:"token.json"`
}

type InvidiousConfig struct {
	Takeout     string `arg:"--invidious.takeout" default:"subscription_manager.json"`
	InstanceURL string `arg:"--invidious.instanceurl,required"`
}

type FreetubeOptions struct {
	Dir string `arg:"--freetube.dir" default:"freetube"`
}

var logger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stderr,
	TimeFormat: "2006/01/02 15:04:05",
}).With().Timestamp().Logger()

func main() {
	cfg := arg.Config{
		Program: "ytmigrator-go",
		Out:     os.Stderr,
	}
	p, err := arg.NewParser(cfg, &args)
	if err != nil {
		panic(err)
	}

	var flags []string
	if len(os.Args) > 0 {
		flags = os.Args[1:]
	}
	err = p.Parse(flags)

	var code int
	//nolint:errcheck
	switch {
	case err == arg.ErrHelp:
		p.WriteHelpForSubcommand(cfg.Out, p.SubcommandNames()...)
	case err != nil:
		p.WriteHelpForSubcommand(cfg.Out, p.SubcommandNames()...)
		os.Stderr.WriteString("error:")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		code = 2
	default: //nolint:staticcheck,gocritic
		p.WriteHelp(cfg.Out)
		code = 2
	case args.Cmd_yt2i != nil:
		code = args.Cmd_yt2i.Run()
	case args.Cmd_yt2ft != nil:
		code = args.Cmd_yt2ft.Run()
	}
	os.Exit(code)
}
