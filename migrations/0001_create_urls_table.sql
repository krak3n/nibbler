-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE urls (
	"id" text PRIMARY KEY,
	"url" text NOT NULL UNIQUE,
	UNIQUE (id, url)
) WITH (OIDS = FALSE);

-- +goose Down
-- SQL in this section is executed when the migration is applied.

DROP TABLE urls;
