// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tubular

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/tubular/internal"
)

func (m *Migrator) ImportTo(data common.UserData) error {
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

func (m *Migrator) importPlaylists(playlists common.Playlists) error {
	for title, videoIDs := range playlists {
		if err := m.importPlaylist(title, videoIDs); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrator) importPlaylist(title common.PlaylistTitle, videoIDs []common.VideoID) error {
	ids := make([]int64, 0, len(videoIDs))
	for _, videoID := range videoIDs {
		url := videoUrl(videoID)
		id, err := m.queries.Stream(context.Background(), url)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v, err := common.Innertube.GetVideo(videoID)
			if err != nil {
				return err
			}
			if v.LiveVideo != nil {
				panic("unimplemented")
			}
			arg := internal.InsertStreamParams{
				Url:         url,
				Title:       v.Video.Title,
				Duration:    v.Video.Duration,
				Uploader:    v.Video.Channel.Name,
				UploaderUrl: nullString(v.Video.Channel.URL),
			}
			if v.Video.ViewCount.IsSome() {
				arg.ViewCount = nullInt64(v.Video.ViewCount.Unwrap())
			}
			uploadTime, err := time.Parse(time.RFC3339, v.Video.UploadDate)
			if err == nil {
				arg.UploadDate = nullInt64(uploadTime.UnixMilli())
			}
			if thumbnail := v.Video.Thumbnails.Best(); thumbnail != nil {
				arg.ThumbnailUrl = nullString(thumbnail.URL)
			}
			id, err = m.queries.InsertStream(context.Background(), arg)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		}
		ids = append(ids, id)
	}

	playlistId, err := m.queries.InsertPlaylist(
		context.Background(),
		internal.InsertPlaylistParams{
			Name:              nullString(title),
			ThumbnailStreamID: ids[0],
		},
	)
	if err != nil {
		return err
	}

	for idx, id := range ids {
		if err = m.queries.InsertPlaylistStreamJoin(
			context.Background(),
			internal.InsertPlaylistStreamJoinParams{
				PlaylistID: playlistId,
				StreamID:   id,
				JoinIndex:  int64(idx),
			},
		); err != nil {
			return err
		}
	}

	return nil
}

func (m Migrator) importSubscriptions(subscriptions common.Subscriptions) error {
	for _, channel := range subscriptions {
		url := nullString(channelUrl(channel.ID))
		has, err := m.queries.HasSubscribed(context.Background(), url)
		if err != nil {
			return err
		}
		if has != 0 {
			continue
		}

		channel, err := common.Innertube.GetChannel(channel.ID)
		if err != nil {
			return err
		}
		arg := internal.InsertSubscriptionParams{
			Url:             url,
			Name:            nullString(channel.Name),
			SubscriberCount: sql.NullInt64{Valid: false},
			Description:     nullString(channel.Description),
		}
		if thumbnail := channel.Thumbnails.Best(); thumbnail != nil {
			arg.AvatarUrl = nullString(thumbnail.URL)
		}
		if err = m.queries.InsertSubscription(context.Background(), arg); err != nil {
			return err
		}
	}
	return nil
}
