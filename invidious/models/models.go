// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package models

import "github.com/go-json-experiment/json/jsontext"

type Takeout struct {
	Subscriptions []string       `json:"subscriptions"`
	WatchHistory  jsontext.Value `json:"watch_history"`
	Preferences   jsontext.Value `json:"preferences"`
	Playlists     []Playlist     `json:"playlists"`
}

type Playlist struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Privacy     string   `json:"privacy"`
	Videos      []string `json:"videos"`
}
