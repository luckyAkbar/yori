-- +migrate Up notransaction

CREATE TABLE IF NOT EXISTS "results" (
    id BIGINT PRIMARY KEY,
    result TEXT DEFAULT ''
);

-- +migrate Down

DROP TABLE IF EXISTS "results";