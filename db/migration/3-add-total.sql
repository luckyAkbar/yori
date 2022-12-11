-- +migrate Up notransaction

ALTER TABLE "results" ADD COLUMN IF NOT EXISTS total_ro BIGINT DEFAULT NULL;
ALTER TABLE "results" ADD COLUMN IF NOT EXISTS total_base BIGINT DEFAULT NULL;

-- +migrate Down

ALTER TABLE "results" DROP COLUMN IF EXISTS total_ro;
ALTER TABLE "results" DROP COLUMN IF EXISTS total_base;