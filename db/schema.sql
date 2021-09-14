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
    -- is_main             BOOLEAN NOT NULL DEFAULT TRUE,
    currency_code       VARCHAR(3) NOT NULL,
    account_name        VARCHAR NOT NULL,
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
//таб для ист транз
CREATE TABLE IF NOT EXISTS transactions (
    id                      bigserial PRIMARY KEY,
    debet_account_id        BIGINT NOT NULL,
    credit_account_id       BIGINT NOT NULL,
    amount                  BIGINT NOT NULL,
    date                    TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
); 
