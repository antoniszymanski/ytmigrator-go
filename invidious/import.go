// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"errors"
	"fmt"
	"sync"

	"github.com/antoniszymanski/invidious-go"
	"github.com/antoniszymanski/ytmigrator-go/common"
	"github.com/antoniszymanski/ytmigrator-go/invidious/models"
	"github.com/cli/browser"
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

func (m *Migrator) importSubscriptions(input common.Subscriptions) error {
	if input == nil {
		m.logger.Info().Msg("importing subscriptions has been omitted")
		return nil
	}

	subscriptions, err := m.client.Subscriptions()
	if e, ok := common.ErrorAs[invidious.Error](err); ok &&
		e.StatusCode == 403 && e.Message == "Endpoint disabled" {
		url := m.client.InstanceURL + "/subscription_manager"
		fmt.Printf(msg, url)
		if err = browser.OpenURL(url); err != nil {
			return err
		}
		common.AwaitEnter() //nolint:errcheck
	} else if err != nil {
		return err
	} else {
		for _, channel := range subscriptions {
			err = m.client.RemoveSubscription(channel.AuthorId)
			if err != nil {
				return err
			}
		}
	}

	m.takeout.Subscriptions = make([]string, 0, len(input))
	for _, channel := range input {
		m.takeout.Subscriptions = append(m.takeout.Subscriptions, channel.ID)
	}
	return nil
}

const msg = `Your selected Invidious instance has disabled "/api/v1/auth/subscriptions" endpoint. You will be redirected to %q. When you remove all your subscriptions, press Enter.\n`

func (m *Migrator) importPlaylists(input common.Playlists) error {
	if input == nil {
		m.logger.Info().Msg("importing playlists has been omitted")
		return nil
	}

	playlists, err := m.client.Playlists()
	if err != nil {
		return err
	}
	for _, playlist := range playlists {
		err = m.client.DeletePlaylist(playlist.PlaylistId)
		if err != nil {
			return err
		}
	}

	m.takeout.Playlists = common.Require(m.takeout.Playlists, len(input))
	for title, videoIDs := range input {
		m.takeout.Playlists = append(m.takeout.Playlists, models.Playlist{
			Title:       title,
			Description: "",
			Privacy:     "Private",
			Videos:      videoIDs,
		})
	}
	return nil
}
