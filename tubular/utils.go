// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tubular

import "database/sql"

func videoUrl(id string) string {
	return "https://www.youtube.com/watch?v=" + id
}

func channelUrl(id string) string {
	return "https://www.youtube.com/channel/" + id
}

func nullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

func nullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: true}
}
