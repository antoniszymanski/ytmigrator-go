// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"sync"

	"github.com/antoniszymanski/innertube-go/youtube"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

var Innertube = sync.OnceValue(func() youtube.Client {
	c, err := youtube.NewClient()
	if err != nil {
		panic(err)
	}
	return c
})
