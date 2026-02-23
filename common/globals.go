// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"sync"

	"github.com/antoniszymanski/innertube-go/youtube"
	"github.com/avast/retry-go/v5"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func GetVideo(id string) (youtube.VideoResult, error) {
	return getVideoRetrier.Do(func() (youtube.VideoResult, error) {
		return client().GetVideo(id)
	})
}

var getVideoRetrier = retry.NewWithData[youtube.VideoResult](
	retry.LastErrorOnly(true),
	retry.Attempts(3),
	retry.RetryIf(func(err error) bool {
		_, ok := err.(youtube.ErrVideoNotFound)
		return !ok
	}),
)

func GetChannel(id string) (*youtube.Channel, error) {
	return getChannelRetrier.Do(func() (*youtube.Channel, error) {
		return client().GetChannel(id)
	})
}

var getChannelRetrier = retry.NewWithData[*youtube.Channel](
	retry.LastErrorOnly(true),
	retry.Attempts(3),
	retry.RetryIf(func(err error) bool {
		_, ok := err.(youtube.ErrChannelNotFound)
		return !ok
	}),
)

var client = sync.OnceValue(func() youtube.Client {
	c, err := youtube.NewClient(nil)
	if err != nil {
		panic(err)
	}
	return c
})
