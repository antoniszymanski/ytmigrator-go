// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tubular

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/antoniszymanski/innertube-go/youtube"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/tubular/internal"
	"github.com/avast/retry-go/v5"
)

func (m *Migrator) ImportTo(data common.UserData) error {
	if err := m.importSubscriptions(data.Subscriptions); err != nil {
		return err
	}
	return m.importPlaylists(data.Playlists)
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

		channel, err := common.Innertube().GetChannel(channel.ID)
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
			v, err := retry.NewWithData[youtube.VideoResult](
				retry.LastErrorOnly(true),
				retry.Attempts(3),
				retry.RetryIf(func(err error) bool {
					_, ok := err.(youtube.ErrVideoNotFound)
					return !ok
				}),
			).Do(func() (youtube.VideoResult, error) {
				return common.Innertube().GetVideo(videoID)
			})
			if _, ok := err.(youtube.ErrVideoNotFound); ok {
				common.Logger.Warn().Str("videoID", videoID).Msg("video not found")
				continue
			}
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
				Uploader:    v.Video.Channel.Unwrap().Name,
				UploaderUrl: nullString(v.Video.Channel.Unwrap().URL),
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
