package youtube

import (
	"context"

	"google.golang.org/api/youtube/v3"
)

// https://developers.google.com/youtube/v3/docs/playlists/list
func (m *Migrator) listPlaylists() ([]*youtube.Playlist, error) {
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

// https://developers.google.com/youtube/v3/docs/playlists/insert
func (m *Migrator) createPlaylist(title string) (*youtube.Playlist, error) {
	return m.client.Playlists.Insert(
		[]string{"snippet", "status"},
		&youtube.Playlist{
			Snippet: &youtube.PlaylistSnippet{
				Title: title,
			},
			Status: &youtube.PlaylistStatus{
				PrivacyStatus: "private",
			},
		},
	).Do()
}

// https://developers.google.com/youtube/v3/docs/playlists/delete
func (m *Migrator) deletePlaylist(playlistID string) error {
	return m.client.Playlists.Delete(playlistID).Do()
}
