DO
$$
    BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
    CREATE DOMAIN UINT256 AS NUMERIC
        CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    ELSE
    ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
    ALTER DOMAIN UINT256 ADD
        CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    END IF;
    END
$$;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" cascade;

CREATE TABLE IF NOT EXISTS exchange(
    guid            TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    name            VARCHAR UNIQUE NOT NULL,
    config          VARCHAR NOT NULL,
    timestamp       INTEGER NOT NULL CHECK (timestamp > 0)
);

CREATE TABLE IF NOT EXISTS support_token(
   guid            TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
   symbol_name     VARCHAR UNIQUE NOT NULL,
   base_asset      VARCHAR NOT NULL,
   qoute_asset     VARCHAR NOT NULL,
   timestamp       INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS support_token_symbol_name ON support_token (symbol_name);

CREATE TABLE IF NOT EXISTS support_token_exchange (
    guid        TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_id    VARCHAR NOT NULL,
    exchange_id VARCHAR NOT NULL
);
CREATE INDEX IF NOT EXISTS support_token_exchange_token_id ON support_token_exchange(token_id);

CREATE TABLE IF NOT EXISTS token_price(
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_id      VARCHAR NOT NULL,
    sell_price    VARCHAR NOT NULL,
    avg_price     VARCHAR NOT NULL,
    rate          VARCHAR NOT NULL,
    volume24h    VARCHAR NOT NULL,
    timestamp     INTEGER NOT NULL CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS token_price_token_id ON token_price(token_id);