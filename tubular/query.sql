-- name: HasSubscribed :one
SELECT
	EXISTS (
		SELECT
			1
		FROM
			"subscriptions"
		WHERE
			"url" = ?
	);

-- name: InsertSubscription :exec
INSERT INTO
	"subscriptions" (
		"service_id",
		"url",
		"name",
		"avatar_url",
		"subscriber_count",
		"description",
		"notification_mode"
	)
VALUES
	(0, ?, ?, ?, ?, ?, 0);

-- name: Stream :one
SELECT
	"uid"
FROM
	"streams"
WHERE
	"url" = ?
LIMIT
	1;

-- name: InsertStream :one
INSERT INTO
	"streams" (
		"service_id",
		"url",
		"title",
		"stream_type",
		"duration",
		"uploader",
		"uploader_url",
		"thumbnail_url",
		"view_count",
		"textual_upload_date",
		"upload_date",
		"is_upload_date_approximation"
	)
VALUES
	(0, ?, ?, "VIDEO_STREAM", ?, ?, ?, ?, ?, ?, ?, 1) RETURNING "uid";

-- name: InsertPlaylist :one
INSERT INTO
	"playlists" (
		"name",
		"is_thumbnail_permanent",
		"thumbnail_stream_id",
		"display_index"
	)
VALUES
	(?, 0, ?, -1) RETURNING "uid";

-- name: InsertPlaylistStreamJoin :exec
INSERT INTO
	"playlist_stream_join" ("playlist_id", "stream_id", "join_index")
VALUES
	(?, ?, ?);