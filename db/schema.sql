-- таблица менеджера
CREATE TABLE IF NOT EXISTS managers (
    id bigserial PRIMARY key,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    phone TEXT NOT NULL unique,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);
-- таблица клиента
CREATE TABLE IF NOT EXISTS customer(
    id              bigserial PRIMARY key,
    name            TEXT NOT NULL,
    surname         TEXT NOT NULL,
    phone           TEXT NOT NULL unique,
    password        TEXT NOT NULL,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created         TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);
-- таблица для счета
CREATE TABLE IF NOT EXISTS account (
    id                  bigserial primary key,
    customer_id         BIGINT NOT NULL REFERENCES customer,
    currency_code       VARCHAR(3),
    account_name        VARCHAR,
    amount              BIGINT
);
-- Таблица для услуги
CREATE TABLE IF NOT EXISTS services (
    id bigserial primary key,
    name VARCHAR NOT NULL
);
-- Таблица для список банкомата
CREATE TABLE IF NOT EXISTS atm (
    id              bigserial primary key,
    numbers         BIGINT NOT NULL,
    district        TEXT NOT NULL,
    address         TEXT NOT NULL
);

DROP TABLE atm;





