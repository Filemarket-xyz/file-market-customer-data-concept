-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TYPE transaction_state AS enum ( 'pending', 'confirmed', 'failed' );

CREATE TABLE public.transactions
(
    id        CHAR(66) PRIMARY KEY,
    state     transaction_state     NOT NULL,
    timestamp BIGINT                NOT NULL
);

CREATE TABLE public.auth_messages
(
    address    VARCHAR(42) NOT NULL,
    code       VARCHAR(255)   NOT NULL,
    created_at BIGINT      NOT NULL,
    CONSTRAINT auth_messages_pkey PRIMARY KEY (address)
);

CREATE TABLE public.users
(
    id       BIGINT     NOT NULL,
    role     INTEGER    NOT NULL,
    address CHAR(42) UNIQUE NOT NULL,
    PRIMARY KEY (id, role)
);

CREATE TABLE public.admins
(
    id  BIGSERIAL PRIMARY KEY
);

CREATE TABLE public.clients
(
    id  BIGSERIAL PRIMARY KEY,
    agreement BOOLEAN NOT NULL,
    bought BOOLEAN NOT NULL,
    point_balance VARCHAR(255) NOT NULL DEFAULT '0'
);

CREATE TABLE public.jwtokens
(
    user_id    BIGINT       NOT NULL,
    role       INTEGER      NOT NULL,
    purpose    INTEGER      NOT NULL,
    number     INTEGER      NOT NULL,
    expires_at BIGINT       NOT NULL,
    secret     varchar(255) NOT NULL,
    FOREIGN KEY (user_id, role) REFERENCES users (id, role)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE public.jwtokens;
DROP TABLE public.clients;
DROP TABLE public.admins;
DROP TABLE public.users;
DROP TABLE public.auth_messages;
DROP TABLE public.transactions;
DROP TYPE transaction_state;