// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package freetube

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/antoniszymanski/stacktrace-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/freetube/models"
	"github.com/go-json-experiment/json"
	"github.com/google/uuid"
	"github.com/kaorimatz/go-opml"
)

func (m *Migrator) ImportTo(data common.UserData) error {
	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(1)
	stacktrace.Go(func() {
		errs[0] = m.importSubscriptions(data.Subscriptions)
		wg.Done()
	}, nil, nil)

	wg.Add(1)
	stacktrace.Go(func() {
		errs[1] = m.importPlaylists(data.Playlists)
		wg.Done()
	}, nil, nil)

	wg.Wait()
	return errors.Join(errs...)
}

func (m *Migrator) importSubscriptions(input common.Subscriptions) error {
	if input == nil {
		common.Logger.Info().Msg("importing subscriptions has been omitted")
		return nil
	}

	output := &opml.OPML{
		Version: "1.1",
		Outlines: []*opml.Outline{{
			Text:     "YouTube Subscriptions",
			Title:    "YouTube Subscriptions",
			Outlines: make([]*opml.Outline, 0, len(input)),
		}},
	}
	outline := output.Outlines[0]
	for _, channel := range input {
		outline.Outlines = append(outline.Outlines, &opml.Outline{
			Text: channel.Name,
			Type: "rss",
			XMLURL: &url.URL{
				Scheme:   "https",
				Host:     "www.youtube.com",
				Path:     "/feeds/videos.xml",
				RawQuery: "channel_id=" + channel.ID,
			},
			Title: channel.Name,
		})
	}

	f, err := os.Create(filepath.Join(m.dir, "subscriptions.opml"))
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	return opml.Render(f, output)
}

func (m *Migrator) importPlaylists(input common.Playlists) error {
	if input == nil {
		common.Logger.Info().Msg("importing playlists has been omitted")
		return nil
	}

	output := make([]models.Playlist, 0, len(input))
	for playlistName, videoIDs := range input {
		playlist := models.Playlist{
			PlaylistName:  playlistName,
			Protected:     false,
			Description:   "",
			Videos:        make([]models.Video, 0, len(videoIDs)),
			ID:            common.Sha256(playlistName),
			CreatedAt:     0,
			LastUpdatedAt: 0,
		}

		for _, videoID := range videoIDs {
			v, err := common.Innertube().GetVideo(videoID)
			if err != nil {
				common.Logger.Warn().Err(err).
					Str("playlistName", playlistName).
					Str("videoID", videoID).
					Msg("failed to add video to the playlist")
				continue
			}
			if v.LiveVideo != nil {
				panic("unimplemented")
			}
			playlist.Videos = append(playlist.Videos, models.Video{
				VideoID:        v.Video.ID,
				Title:          v.Video.Title,
				Author:         v.Video.Channel.Name,
				AuthorID:       v.Video.Channel.ID,
				LengthSeconds:  v.Video.Duration,
				TimeAdded:      0,
				PlaylistItemID: uuid.New().String(),
				Type:           "video",
			})
			common.Logger.Debug().
				Str("playlistName", playlistName).
				Str("videoID", videoID).
				Msg("added video to the playlist")
		}
		output = append(output, playlist)
	}

	f, err := os.Create(filepath.Join(m.dir, "playlists.db"))
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	for _, playlist := range output {
		if err = json.MarshalWrite(f, &playlist); err != nil {
			return err
		}
		if _, err = f.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}
