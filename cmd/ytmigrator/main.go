// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/alexflint/go-arg"
	"github.com/antoniszymanski/stacktrace-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/rs/zerolog"
)

type Cli struct {
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

type LevelOption zerolog.Level

var _ encoding.TextUnmarshaler = (*LevelOption)(nil)

func (opt *LevelOption) UnmarshalText(text []byte) (err error) {
	switch got := string(text); got {
	case "debug":
		*opt = LevelOption(zerolog.DebugLevel)
	case "info":
		*opt = LevelOption(zerolog.InfoLevel)
	case "warn":
		*opt = LevelOption(zerolog.WarnLevel)
	case "error":
		*opt = LevelOption(zerolog.ErrorLevel)
	default:
		err = notOneOf(got, "debug", "info", "warn", "error")
	}
	return err
}

type TypeOption struct {
	Writer io.Writer
}

var _ encoding.TextUnmarshaler = (*TypeOption)(nil)

func (opt *TypeOption) UnmarshalText(text []byte) (err error) {
	switch got := string(text); got {
	case "json":
		opt.Writer = os.Stderr
	case "console":
		opt.Writer = zerolog.ConsoleWriter{
			Out:         os.Stderr,
			TimeFormat:  "2006/01/02 15:04:05",
			FieldsOrder: []string{"error"},
		}
	default:
		err = notOneOf(got, "json", "console")
	}
	return err
}

func notOneOf(got string, members ...string) error {
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
		Program: "ytmigrator",
		Out:     os.Stderr,
	}
	var cli Cli
	p, err := arg.NewParser(cfg, &cli)
	if err != nil {
		panic(err)
	}
	var flags []string
	if len(os.Args) > 0 {
		flags = os.Args[1:]
	}
	err = p.Parse(flags)
	if err == nil {
		common.Logger = zerolog.New(cli.Type.Writer).With().Timestamp().Logger()
		common.Logger = common.Logger.Level(zerolog.Level(cli.Level))
	}
	var code int
	switch {
	case err == arg.ErrHelp:
		p.WriteHelp(cfg.Out)
	case err != nil:
		p.WriteHelp(cfg.Out)
		printErr(err)
		code = 2
	default: //nolint:staticcheck,gocritic
		p.WriteHelp(cfg.Out)
		code = 2
	case cli.Cmd_version != nil:
		code = cli.Cmd_version.Run()
	case cli.Cmd_yt2i != nil:
		code = cli.Cmd_yt2i.Run()
	case cli.Cmd_yt2ft != nil:
		code = cli.Cmd_yt2ft.Run()
	case cli.Cmd_yt2t != nil:
		code = cli.Cmd_yt2t.Run()
	}
	os.Exit(code) //nolint:gocritic // exitAfterDefer
}

//nolint:errcheck
func printErr(err error) {
	os.Stderr.WriteString(err.Error())
	os.Stderr.WriteString("\n")
}
