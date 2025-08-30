// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package youtube

import (
	"errors"
	"sync"

	"github.com/antoniszymanski/stacktrace-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
)

func (m *Migrator) ExportFrom(opts common.ExportOptions) (common.UserData, error) {
	var data common.UserData
	var wg sync.WaitGroup
	errs := make([]error, 2)
	if !opts.SkipSubscriptions {
		wg.Add(1)
		stacktrace.Go(func() {
			data.Subscriptions, errs[0] = m.exportSubscriptions()
			wg.Done()
		}, nil, nil)
	}
	if !opts.SkipPlaylists {
		wg.Add(1)
		stacktrace.Go(func() {
			data.Playlists, errs[1] = m.exportPlaylists()
			wg.Done()
		}, nil, nil)
	}
	wg.Wait()
	return data, errors.Join(errs...)
}

func (m *Migrator) exportSubscriptions() (common.Subscriptions, error) {
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

func (m *Migrator) exportPlaylists() (common.Playlists, error) {
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
