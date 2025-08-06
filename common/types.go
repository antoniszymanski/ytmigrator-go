// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/rs/zerolog"

type Migrator interface {
	Import(data UserData) error
	Export(opts ExportOptions) (data UserData, err error)
	SetLogger(logger *zerolog.Logger)
}

type ExportOptions struct {
	SkipSubscriptions bool `arg:"--skip-subscriptions"`
	SkipPlaylists     bool `arg:"--skip-playlists"`
}

type UserData struct {
	Subscriptions Subscriptions `json:"subscriptions"`
	Playlists     Playlists     `json:"playlists"`
}

type Subscriptions []Channel

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Playlists map[PlaylistTitle][]VideoID

type PlaylistTitle = string

type VideoID = string
