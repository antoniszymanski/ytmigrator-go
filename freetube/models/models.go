// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package models

type Playlists []Playlist

type Playlist struct {
	PlaylistName  string  `json:"playlistName"`
	Protected     bool    `json:"protected"`
	Description   string  `json:"description"`
	Videos        []Video `json:"videos"`
	ID            string  `json:"_id"`
	CreatedAt     int64   `json:"createdAt"`
	LastUpdatedAt int64   `json:"lastUpdatedAt"`
}

type Video struct {
	VideoID        string `json:"videoId"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	AuthorID       string `json:"authorId"`
	LengthSeconds  int    `json:"lengthSeconds"`
	TimeAdded      int64  `json:"timeAdded"`
	PlaylistItemID string `json:"playlistItemId"`
	Type           string `json:"type"`
}

type Subscriptions struct {
	ID            string                `json:"_id"`
	Name          string                `json:"name"`
	BgColor       string                `json:"bgColor"`
	TextColor     string                `json:"textColor"`
	Subscriptions []SubscriptionChannel `json:"subscriptions"`
}

type SubscriptionChannel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}
