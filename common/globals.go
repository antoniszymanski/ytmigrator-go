// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/antoniszymanski/innertube-go/youtube"

var Innertube = Must(youtube.NewClient())

func Must[A any](a A, err error) A {
	if err != nil {
		panic(err)
	}
	return a
}
