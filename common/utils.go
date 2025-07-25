// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"unsafe"

	"github.com/raitonoberu/ytsearch"
)

func ChannelSearch(query string) ([]*ytsearch.ChannelItem, error) {
	search := ytsearch.ChannelSearch(query)
	result, err := search.Next()
	if err != nil {
		return nil, err
	}
	return result.Channels, nil
}

func PlaylistSearch(query string) ([]*ytsearch.PlaylistItem, error) {
	search := ytsearch.PlaylistSearch(query)
	result, err := search.Next()
	if err != nil {
		return nil, err
	}
	return result.Playlists, nil
}

func VideoSearch(query string) ([]*ytsearch.VideoItem, error) {
	search := ytsearch.VideoSearch(query)
	result, err := search.Next()
	if err != nil {
		return nil, err
	}
	return result.Videos, nil
}

func ErrorAs[T error](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}

func Require[T any](s []T, length int) []T {
	if length <= cap(s) {
		return s[:0]
	}
	return make([]T, 0, length)
}

func Sha256(s string) string {
	src := sha256.Sum256(string2bytes(s))
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src[:])
	return bytes2string(dst)
}

func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func bytes2string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
