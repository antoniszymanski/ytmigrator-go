// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"unsafe"
)

func Resize[S ~[]E, E any](s S, length int) S {
	if length <= cap(s) {
		return s[:length]
	}
	return append(s[:cap(s)], make(S, length-cap(s))...)
}

func Require[S ~[]E, E any](s S, capacity int) S {
	if capacity <= cap(s) {
		return s[:0]
	}
	return make(S, 0, capacity)
}

func ErrorAs[T error](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}

func Sha256(s string) string {
	src := sha256.Sum256(StringToBytes(s))
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src[:])
	return BytesToString(dst)
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
