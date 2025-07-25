// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package youtube

import (
	"context"
	"errors"
	"sync"

	"github.com/antoniszymanski/ytmigrator-go/common"
	"google.golang.org/api/youtube/v3"
)

func (m Migrator) Export(opts common.ExportOptions) (common.UserData, error) {
	var data common.UserData
	var wg sync.WaitGroup
	errs := make([]error, 2)
	if !opts.SkipSubscriptions {
		wg.Add(1)
		go func() {
			data.Subscriptions, errs[0] = m.exportSubscriptions()
			wg.Done()
		}()
	}
	if !opts.SkipPlaylists {
		wg.Add(1)
		go func() {
			data.Playlists, errs[1] = m.exportPlaylists()
			wg.Done()
		}()
	}
	wg.Wait()
	return data, errors.Join(errs...)
}

func (m Migrator) exportSubscriptions() (common.Subscriptions, error) {
	subscriptions, err := m.listSubscriptions()
	if err != nil {
		return nil, err
	}

	output := make(common.Subscriptions, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		output = append(output, common.Channel{
			ID:   subscription.Snippet.ResourceId.ChannelId,
			Name: subscription.Snippet.Title,
		})
	}
	return output, nil
}

func (m Migrator) listSubscriptions() ([]*youtube.Subscription, error) {
	var items []*youtube.Subscription
	f := func(resp *youtube.SubscriptionListResponse) error {
		items = append(items, resp.Items...)
		return nil
	}

	err := m.client.Subscriptions.List([]string{"snippet"}).
		MaxResults(50).
		Mine(true).
		Pages(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (m Migrator) exportPlaylists() (common.Playlists, error) {
	playlists, err := m.listPlaylists()
	if err != nil {
		return nil, err
	}

	output := make(common.Playlists, len(playlists))
	for _, playlist := range playlists {
		items, err := m.listPlaylistItems(playlist.Id)
		if err != nil {
			return nil, err
		}

		videoIDs := make([]common.VideoID, 0, len(items))
		for _, item := range items {
			videoIDs = append(videoIDs, item.Snippet.ResourceId.VideoId)
		}
		output[playlist.Snippet.Title] = videoIDs
	}

	return output, nil
}

func (m Migrator) listPlaylists() ([]*youtube.Playlist, error) {
	var items []*youtube.Playlist
	f := func(resp *youtube.PlaylistListResponse) error {
		items = append(items, resp.Items...)
		return nil
	}

	err := m.client.Playlists.List([]string{"snippet"}).
		MaxResults(50).
		Mine(true).
		Pages(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (m Migrator) listPlaylistItems(playlistId string) ([]*youtube.PlaylistItem, error) {
	var items []*youtube.PlaylistItem
	f := func(resp *youtube.PlaylistItemListResponse) error {
		items = append(items, resp.Items...)
		return nil
	}

	err := m.client.PlaylistItems.List([]string{"snippet"}).
		MaxResults(50).
		PlaylistId(playlistId).
		Pages(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return items, nil
}
