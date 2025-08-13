package youtube

import (
	"context"

	"google.golang.org/api/youtube/v3"
)

// https://developers.google.com/youtube/v3/docs/playlistItems/list
func (m *Migrator) listPlaylistItems(playlistID string) ([]*youtube.PlaylistItem, error) {
	var items []*youtube.PlaylistItem
	f := func(resp *youtube.PlaylistItemListResponse) error {
		items = append(items, resp.Items...)
		return nil
	}

	err := m.client.PlaylistItems.List([]string{"snippet"}).
		MaxResults(50).
		PlaylistId(playlistID).
		Pages(context.Background(), f)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// https://developers.google.com/youtube/v3/docs/playlistItems/insert
func (m *Migrator) insertPlaylistItem(playlistID, videoID string, position int64) (*youtube.PlaylistItem, error) {
	item := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}
	if position >= 0 {
		item.Snippet.Position = position
		item.Snippet.ForceSendFields = []string{"Position"}
	}
	return m.client.PlaylistItems.Insert([]string{"snippet"}, item).Do()
}

// https://developers.google.com/youtube/v3/docs/playlistItems/update
func (m *Migrator) updatePlaylistItem(
	playlistID,
	itemID,
	videoID string,
	position int64,
) (*youtube.PlaylistItem, error) {
	item := &youtube.PlaylistItem{
		Id: itemID,
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}
	if position >= 0 {
		item.Snippet.Position = position
		item.Snippet.ForceSendFields = []string{"Position"}
	}
	item, err := m.client.PlaylistItems.Update([]string{"id", "snippet"}, item).Do()
	if err == nil {
		return item, nil
	}

	if err = m.deletePlaylistItem(itemID); err != nil {
		return nil, err
	}
	return m.insertPlaylistItem(playlistID, videoID, position)
}

// https://developers.google.com/youtube/v3/docs/playlistItems/delete
func (m *Migrator) deletePlaylistItem(id string) error {
	return m.client.PlaylistItems.Delete(id).Do()
}
