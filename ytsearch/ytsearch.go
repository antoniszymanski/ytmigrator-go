package ytsearch

import (
	"errors"
	"slices"

	"github.com/raitonoberu/ytsearch"
)

func FindChannelByID(id string) (*ytsearch.ChannelItem, error) {
	search := ytsearch.ChannelSearch(id)
	for search.NextExists() {
		result, err := search.Next()
		if err != nil {
			return nil, err
		}
		i := slices.IndexFunc(result.Channels,
			func(item *ytsearch.ChannelItem) bool {
				return item.ID == id
			},
		)
		if i != -1 {
			return result.Channels[i], nil
		}
	}
	return nil, ErrNotFound
}

func FindPlaylistByID(id string) (*ytsearch.PlaylistItem, error) {
	search := ytsearch.PlaylistSearch(id)
	for search.NextExists() {
		result, err := search.Next()
		if err != nil {
			return nil, err
		}
		i := slices.IndexFunc(result.Playlists,
			func(item *ytsearch.PlaylistItem) bool {
				return item.ID == id
			},
		)
		if i != -1 {
			return result.Playlists[i], nil
		}
	}
	return nil, ErrNotFound
}

func FindVideoByID(id string) (*ytsearch.VideoItem, error) {
	search := ytsearch.VideoSearch(id)
	for search.NextExists() {
		result, err := search.Next()
		if err != nil {
			return nil, err
		}
		i := slices.IndexFunc(result.Videos,
			func(item *ytsearch.VideoItem) bool {
				return item.ID == id
			},
		)
		if i != -1 {
			return result.Videos[i], nil
		}
	}
	return nil, ErrNotFound
}

var ErrNotFound = errors.New("not found")
