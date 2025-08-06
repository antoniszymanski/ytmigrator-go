// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"unsafe"
)

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
	return Bytes2string(dst)
}

func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Bytes2string(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
