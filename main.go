// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding"
	"errors"
	"io"
	"os"
	"slices"
	"strconv"

	"github.com/alexflint/go-arg"
	"github.com/antoniszymanski/stacktrace-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/rs/zerolog"
)

var args struct {
	LoggingOptions
	*Cmd_version `arg:"subcommand:version" help:"display version and exit"`
	*Cmd_yt2i    `arg:"subcommand:yt2i"`
	*Cmd_yt2ft   `arg:"subcommand:yt2ft"`
	*Cmd_yt2t    `arg:"subcommand:yt2t"`
}

type LoggingOptions struct {
	Level LevelOption `arg:"--logging.level" default:"info"`
	Type  TypeOption  `arg:"--logging.type" default:"console"`
}

type LevelOption string

var _ encoding.TextUnmarshaler = (*LevelOption)(nil)

func (opt *LevelOption) UnmarshalText(text []byte) error {
	*opt = LevelOption(text)
	return validateOneOf(string(text), "debug", "info", "warn", "error")
}

type TypeOption string

var _ encoding.TextUnmarshaler = (*TypeOption)(nil)

func (opt *TypeOption) UnmarshalText(text []byte) error {
	*opt = TypeOption(text)
	return validateOneOf(string(text), "json", "console")
}

func validateOneOf(got string, members ...string) error {
	if slices.Contains(members, got) {
		return nil
	}
	var dst []byte
	dst = append(dst, "must be one of "...)
	for i, member := range members {
		dst = strconv.AppendQuote(dst, member)
		if i < len(members)-1 {
			dst = append(dst, ',')
		}
	}
	dst = append(dst, " but got "...)
	dst = strconv.AppendQuote(dst, got)
	return errors.New(common.BytesToString(dst))
}

type YoutubeOptions struct {
	Credentials string `arg:"--youtube.credentials" default:"credentials.json"`
	Token       string `arg:"--youtube.token" default:"token.json"`
}

type InvidiousOptions struct {
	Takeout     string `arg:"--invidious.takeout" default:"subscription_manager.json"`
	InstanceURL string `arg:"--invidious.instanceurl,required"`
}

type FreetubeOptions struct {
	Dir string `arg:"--freetube.dir" default:"freetube"`
}

type TubularOptions struct {
	DSN string `arg:"--tubular.dsn,required"`
}

func main() {
	defer stacktrace.Handle(true, nil, nil)

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

	if err == nil {
		var w io.Writer
		switch args.Type {
		case "json":
			w = os.Stderr
		case "console":
			w = zerolog.ConsoleWriter{
				Out:         os.Stderr,
				TimeFormat:  "2006/01/02 15:04:05",
				FieldsOrder: []string{"error"},
			}
		}
		common.Logger = zerolog.New(w).With().Timestamp().Logger()
		switch args.Level {
		case "debug":
			common.Logger = common.Logger.Level(zerolog.DebugLevel)
		case "info":
			common.Logger = common.Logger.Level(zerolog.InfoLevel)
		case "warn":
			common.Logger = common.Logger.Level(zerolog.WarnLevel)
		case "error":
			common.Logger = common.Logger.Level(zerolog.ErrorLevel)
		}
	}

	var code int
	//nolint:errcheck
	switch {
	case err == arg.ErrHelp:
		p.WriteHelpForSubcommand(cfg.Out, p.SubcommandNames()...)
	case err != nil:
		p.WriteHelpForSubcommand(cfg.Out, p.SubcommandNames()...)
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		code = 2
	default: //nolint:staticcheck,gocritic
		p.WriteHelp(cfg.Out)
		code = 2
	case args.Cmd_version != nil:
		code = args.Cmd_version.Run()
	case args.Cmd_yt2i != nil:
		code = args.Cmd_yt2i.Run()
	case args.Cmd_yt2ft != nil:
		code = args.Cmd_yt2ft.Run()
	case args.Cmd_yt2t != nil:
		code = args.Cmd_yt2t.Run()
	}
	os.Exit(code) //nolint:gocritic // exitAfterDefer
}
