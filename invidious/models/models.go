// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package models

import "github.com/go-json-experiment/json/jsontext"

type Takeout struct {
	Subscriptions []string       `json:"subscriptions"`
	Playlists     []Playlist     `json:"playlists"`
	Unknown       jsontext.Value `json:",unknown"`
}

type Playlist struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Privacy     string   `json:"privacy"`
	Videos      []string `json:"videos"`
}
