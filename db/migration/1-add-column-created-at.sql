-- +migrate Up notransaction
CREATE SEQUENCE record_id_seq START 1 INCREMENT 1;

CREATE TABLE IF NOT EXISTS "new_records" (
    id SERIAL PRIMARY KEY,
    no BIGINT NOT NULL,
    nama TEXT DEFAULT NULL,
    no_engine TEXT DEFAULT NULL,
    tgl_mohon_faktur TIMESTAMP DEFAULT NULL,
    fincoy TEXT DEFAULT NULL,
    type TEXT DEFAULT NULL,
    ktp TEXT DEFAULT NULL,
    kk TEXT DEFAULT NULL,
    dealer TEXT DEFAULT NULL,
    bulan INT DEFAULT NULL,
    tahun INT DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS "new_records";