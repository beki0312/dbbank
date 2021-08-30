CREATE TABLE IF NOT EXISTS customer(
    id bigserial PRIMARY key,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    phone TEXT NOT NULL unique,
    password TEXT NOT NULL,
    amount INTEGER,    
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS account (
    id             bigserial primary key,
    account        bigint,
    amount         BIGINT
);
CREATE TABLE IF NOT EXISTS services (
    id bigserial primary key,
    name VARCHAR(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS atm (
    id              bigserial primary key,
    numbers         BIGINT NOT NULL,
    district        TEXT NOT NULL,
    address         TEXT NOT NULL,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
DROP TABLE atm;
CREATE TABLE IF NOT EXISTS managers (
    id bigserial PRIMARY key,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    phone TEXT NOT NULL unique,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);




