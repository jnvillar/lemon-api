create table user
(
    id           varchar(36)  not null,
    firstname    varchar(255) not null,
    lastname     varchar(255) not null,
    alias        varchar(255) not null unique,
    email        varchar(255) not null unique,
    date_created timestamp,
    primary key (id)
);

create table wallet
(
    id       varchar(36)  not null,
    user_id  varchar(36)  not null,
    currency varchar(255) not null,
    balance  int          not null CHECK (balance >= 0),
    primary key (id)
);

CREATE INDEX wallet_user_id ON wallet (user_id);

create table transaction
(
    id           varchar(36)  not null,
    currency     varchar(255) not null,
    amount       int          not null CHECK (amount >= 0),
    type         varchar(255) not null,
    user_id      varchar(36)  not null,
    user_from    varchar(36)  not null,
    user_to      varchar(36)  not null,
    wallet_id   varchar(36)  not null,
    wallet_to   varchar(36)  not null,
    wallet_from varchar(36)  not null,
    date_created timestamp,
    primary key (id)
);

CREATE INDEX transaction_user_id ON transaction (user_id);
CREATE INDEX transaction_wallet_id ON transaction (wallet_id);