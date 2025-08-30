// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

type Migrator interface {
	ImportTo(data UserData) error
	ExportFrom(opts ExportOptions) (data UserData, err error)
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
