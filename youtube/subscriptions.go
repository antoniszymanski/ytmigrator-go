package youtube

import (
	"context"

	"google.golang.org/api/youtube/v3"
)

// https://developers.google.com/youtube/v3/docs/subscriptions/list
func (m *Migrator) listSubscriptions() ([]*youtube.Subscription, error) {
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

// https://developers.google.com/youtube/v3/docs/subscriptions/insert
func (m *Migrator) insertSubscription(channelID string) error {
	_, err := m.client.Subscriptions.Insert(
		[]string{"snippet"},
		&youtube.Subscription{
			Snippet: &youtube.SubscriptionSnippet{
				ResourceId: &youtube.ResourceId{
					ChannelId: channelID,
				},
			},
		},
	).Do()
	return err
}

// https://developers.google.com/youtube/v3/docs/subscriptions/delete
func (m *Migrator) deleteSubscription(id string) error {
	return m.client.Subscriptions.Delete(id).Do()
}
