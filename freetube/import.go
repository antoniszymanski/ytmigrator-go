// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package freetube

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/antoniszymanski/innertube-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/freetube/models"
	"github.com/go-json-experiment/json"
	"github.com/google/uuid"
	"github.com/kaorimatz/go-opml"
	"golang.org/x/sync/errgroup"
)

func (m Migrator) Import(data common.UserData) error {
	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(1)
	go func() {
		errs[0] = m.importSubscriptions(data.Subscriptions)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errs[1] = m.importPlaylists(data.Playlists)
		wg.Done()
	}()

	wg.Wait()
	return errors.Join(errs...)
}

func (m Migrator) importSubscriptions(input common.Subscriptions) error {
	if input == nil {
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
		return nil
	}

	output := make([]models.Playlist, 0, len(input))
	var mu sync.Mutex

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(30)

	for playlistName, videoIDs := range input {
		g.Go(func() error {
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
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					var c innertube.Client
					v, err := c.GetVideo(videoID)
					if err != nil {
						return err
					}
					playlist.Videos = append(playlist.Videos, models.Video{
						VideoID:        v.ID,
						Title:          v.Title,
						Author:         v.Author,
						AuthorID:       v.AuthorID,
						LengthSeconds:  int(v.LengthSeconds),
						TimeAdded:      0,
						PlaylistItemID: uuid.New().String(),
						Type:           "video",
					})
				}
			}

			mu.Lock()
			output = append(output, playlist)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
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
