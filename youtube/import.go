// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package youtube

import (
	"errors"
	"slices"
	"strconv"
	"sync"

	"github.com/antoniszymanski/stacktrace-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/r3labs/diff/v3"
	"google.golang.org/api/youtube/v3"
)

func (m *Migrator) ImportTo(data common.UserData) error {
	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(1)
	stacktrace.Go(func() {
		errs[0] = m.importSubscriptions(data.Subscriptions)
		wg.Done()
	}, nil)

	wg.Add(1)
	stacktrace.Go(func() {
		errs[1] = m.importPlaylists(data.Playlists)
		wg.Done()
	}, nil)

	wg.Wait()
	return errors.Join(errs...)
}

func (m *Migrator) importSubscriptions(input common.Subscriptions) error {
	if input == nil {
		m.logger.Info().Msg("importing subscriptions has been omitted")
		return nil
	}

	subscriptions, err := m.listSubscriptions()
	if err != nil {
		return err
	}
	src := make([]string, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		src = append(src, subscription.Snippet.ResourceId.ChannelId)
	}

	dst := make([]string, 0, len(input))
	for _, channel := range input {
		dst = append(dst, channel.ID)
	}

	var d diff.Differ
	cl, err := d.Diff(src, dst)
	if err != nil {
		return err
	}
	for _, c := range cl {
		i, err := strconv.Atoi(c.Path[0])
		if err != nil {
			return err
		}
		switch c.Type {
		case diff.CREATE:
			err = m.insertSubscription(c.To.(string))
		case diff.DELETE:
			err = m.deleteSubscription(subscriptions[i].Id)
		default:
			panic("unreachable")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) importPlaylists(input common.Playlists) error {
	if input == nil {
		m.logger.Info().Msg("importing playlists has been omitted")
		return nil
	}

	playlists, err := m.listPlaylists()
	if err != nil {
		return err
	}
	type entry struct {
		Playlist *youtube.Playlist
		Items    []*youtube.PlaylistItem
	}
	entries := make(map[common.PlaylistTitle]entry, len(playlists))
	for _, playlist := range playlists {
		items, err := m.listPlaylistItems(playlist.Id)
		if err != nil {
			return err
		}
		entries[playlist.Snippet.Title] = entry{
			Playlist: playlist,
			Items:    items,
		}
	}

	src := make(common.Playlists, len(entries))
	for title, entry := range entries {
		videoIDs := make([]common.VideoID, 0, len(entry.Items))
		for _, item := range entry.Items {
			videoIDs = append(videoIDs, item.Snippet.ResourceId.VideoId)
		}
		src[title] = videoIDs
	}

	d := &diff.Differ{SliceOrdering: true}
	cl, err := d.Diff(src, input)
	if err != nil {
		return err
	}

	createPlaylist := func(c diff.Change) error {
		var entry entry
		var err error
		entry.Playlist, err = m.createPlaylist(c.Path[0])
		if err != nil {
			return err
		}
		for i, videoID := range c.To.([]common.VideoID) {
			item, err := m.insertPlaylistItem(entry.Playlist.Id, videoID, int64(i))
			if err != nil {
				return err
			}
			entry.Items = append(entry.Items, item)
		}
		entries[entry.Playlist.Snippet.Title] = entry
		return nil
	}
	deletePlaylist := func(c diff.Change) error {
		err := m.deletePlaylist(entries[c.Path[0]].Playlist.Id)
		if err == nil {
			delete(entries, c.Path[0])
		}
		return err
	}
	createPlaylistItem := func(c diff.Change) error {
		entry := entries[c.Path[0]]
		pos, err := strconv.ParseInt(c.Path[1], 10, 64)
		if err != nil {
			return err
		}
		item, err := m.insertPlaylistItem(entry.Playlist.Id, c.To.(string), pos)
		if err != nil {
			return err
		}
		entry.Items = common.Resize(entry.Items, int(pos+1))
		entry.Items[pos] = item
		entries[c.Path[0]] = entry
		return nil
	}
	updatePlaylistItem := func(c diff.Change) error {
		entry := entries[c.Path[0]]
		pos, err := strconv.ParseInt(c.Path[1], 10, 64)
		if err != nil {
			return err
		}
		item, err := m.updatePlaylistItem(
			entry.Playlist.Id,
			entry.Items[pos].Id,
			c.To.(string),
			pos,
		)
		if err != nil {
			return err
		}
		entry.Items[pos] = item
		entries[c.Path[0]] = entry
		return nil
	}
	deletePlaylistItem := func(c diff.Change) error {
		entry := entries[c.Path[0]]
		i, err := strconv.Atoi(c.Path[1])
		if err != nil {
			return err
		}
		if err = m.deletePlaylistItem(entry.Items[i].Id); err != nil {
			return err
		}
		entry.Items = slices.Delete(entry.Items, i, i+1)
		entries[c.Path[0]] = entry
		return nil
	}

	for _, c := range cl {
		var err error
		switch len(c.Path) {
		case 1:
			switch c.Type {
			case diff.CREATE:
				err = createPlaylist(c)
			case diff.DELETE:
				err = deletePlaylist(c)
			default:
				panic("unreachable")
			}
		case 2:
			switch c.Type {
			case diff.CREATE:
				err = createPlaylistItem(c)
			case diff.UPDATE:
				err = updatePlaylistItem(c)
			case diff.DELETE:
				err = deletePlaylistItem(c)
			default:
				panic("unreachable")
			}
		default:
			panic("unreachable")
		}
		if err != nil {
			return err
		}
	}

	return nil
}
