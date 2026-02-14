CREATE TABLE "android_metadata" ("locale" TEXT);

CREATE TABLE "feed" (
	"stream_id" INTEGER NOT NULL,
	"subscription_id" INTEGER NOT NULL,
	PRIMARY KEY ("stream_id", "subscription_id"),
	FOREIGN KEY ("stream_id") REFERENCES "streams" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "feed_group" (
	"uid" INTEGER NOT NULL,
	"name" TEXT NOT NULL,
	"icon_id" INTEGER NOT NULL,
	"sort_order" INTEGER NOT NULL
);

CREATE TABLE "feed_group_subscription_join" (
	"group_id" INTEGER NOT NULL,
	"subscription_id" INTEGER NOT NULL,
	PRIMARY KEY ("group_id", "subscription_id"),
	FOREIGN KEY ("group_id") REFERENCES "feed_group" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "feed_last_updated" (
	"subscription_id" INTEGER NOT NULL,
	"last_updated" INTEGER,
	PRIMARY KEY ("subscription_id"),
	FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "playlist_stream_join" (
	"playlist_id" INTEGER NOT NULL,
	"stream_id" INTEGER NOT NULL,
	"join_index" INTEGER NOT NULL,
	PRIMARY KEY ("playlist_id", "join_index"),
	FOREIGN KEY ("playlist_id") REFERENCES "playlists" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
	FOREIGN KEY ("stream_id") REFERENCES "streams" ("uid") ON UPDATE CASCADE ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "playlists" (
	"uid" INTEGER NOT NULL,
	"name" TEXT,
	"is_thumbnail_permanent" INTEGER NOT NULL,
	"thumbnail_stream_id" INTEGER NOT NULL,
	"display_index" INTEGER NOT NULL
);

CREATE TABLE "remote_playlists" (
	"uid" INTEGER NOT NULL,
	"service_id" INTEGER NOT NULL,
	"name" TEXT,
	"url" TEXT,
	"thumbnail_url" TEXT,
	"uploader" TEXT,
	"display_index" INTEGER NOT NULL,
	"stream_count" INTEGER
);

CREATE TABLE "room_master_table" (
	"id" INTEGER,
	"identity_hash" TEXT,
	PRIMARY KEY ("id")
);

CREATE TABLE "search_history" (
	"creation_date" INTEGER,
	"service_id" INTEGER NOT NULL,
	"search" TEXT,
	"id" INTEGER NOT NULL
);

CREATE TABLE "stream_history" (
	"stream_id" INTEGER NOT NULL,
	"access_date" INTEGER NOT NULL,
	"repeat_count" INTEGER NOT NULL,
	PRIMARY KEY ("stream_id", "access_date"),
	FOREIGN KEY ("stream_id") REFERENCES "streams" ("uid") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "stream_state" (
	"stream_id" INTEGER NOT NULL,
	"progress_time" INTEGER NOT NULL,
	PRIMARY KEY ("stream_id"),
	FOREIGN KEY ("stream_id") REFERENCES "streams" ("uid") ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "streams" (
	"uid" INTEGER NOT NULL,
	"service_id" INTEGER NOT NULL,
	"url" TEXT NOT NULL,
	"title" TEXT NOT NULL,
	"stream_type" TEXT NOT NULL,
	"duration" INTEGER NOT NULL,
	"uploader" TEXT NOT NULL,
	"uploader_url" TEXT,
	"thumbnail_url" TEXT,
	"view_count" INTEGER,
	"textual_upload_date" TEXT,
	"upload_date" INTEGER,
	"is_upload_date_approximation" INTEGER
);

CREATE TABLE "subscriptions" (
	"uid" INTEGER NOT NULL,
	"service_id" INTEGER NOT NULL,
	"url" TEXT,
	"name" TEXT,
	"avatar_url" TEXT,
	"subscriber_count" INTEGER,
	"description" TEXT,
	"notification_mode" INTEGER NOT NULL
);