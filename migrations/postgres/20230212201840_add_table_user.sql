-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id                  serial PRIMARY KEY,
    telegram_id         bigint NOT NULL UNIQUE,
    telegram_chat_id    bigint NOT NULL,
    telegram_first_name varchar,
    telegram_last_name  varchar,
    telegram_user_name  varchar,
    created_at          timestamp    DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS games
(
    id          serial PRIMARY KEY,
    name        varchar NOT NULL,
    url         varchar NOT NULL UNIQUE,
    created_at  timestamp    DEFAULT NOW(),
    deleted_at  timestamp    DEFAULT NULL
);

CREATE TYPE subscription_type AS ENUM ('usual', 'ps_plus', 'ea');

CREATE TABLE IF NOT EXISTS prices
(
    id          serial PRIMARY KEY,
    game_id     bigint NOT NULL
        REFERENCES games(id),
    value       real NOT NULL,
    is_free     bool NOT NULL DEFAULT false,
    type        subscription_type NOT NULL,
    currency    varchar NOT NULL,
    created_at  timestamp    DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users_games
(
    user_telegram_id    integer NOT NULL
        REFERENCES users(telegram_id),
    game_id             integer NOT NULL
        REFERENCES games(id),
    subscription_price  real NOT NULL,
    created_at          timestamp    DEFAULT NOW(),
    deleted_at          timestamp    DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_games;
DROP TABLE IF EXISTS prices;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS games;
DROP TYPE subscription_type;
-- +goose StatementEnd
