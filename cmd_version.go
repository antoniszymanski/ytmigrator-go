package main

import (
	"os"
	"runtime"
	"runtime/debug"
)

type Cmd_version struct{}

func (Cmd_version) Run() int {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		os.Stderr.WriteString("error: build info not found\n") //nolint:errcheck
		return 1
	}

	var revision, time string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			time = setting.Value
		}
	}
	if revision == "" {
		revision = "unknown"
	}
	if time == "" {
		time = "unknown"
	}

	s := "\n" +
		"  Version          " + info.Main.Version + "\n" +
		"  Git Commit       " + revision + "\n" +
		"  Commit Date      " + time + "\n" +
		"  Go Version       " + info.GoVersion + "\n" +
		"  Platform         " + runtime.GOOS + "/" + runtime.GOARCH + "\n" +
		"\n"
	os.Stdout.WriteString(s) //nolint:errcheck

	return 0
}
