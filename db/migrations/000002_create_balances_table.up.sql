create table if not exists balances(
    id bigserial primary key,
    user_id bigint not null references users(id),
    balance numeric(20, 2) not null default 0,
    currency varchar(3) not null default 'USD',
    created_at timestamptz not null default current_timestamp,
    updated_at timestamptz not null default current_timestamp
);

-- add index to balances table
create index on balances(user_id);

-- add unique index composed of user_id and currency
create unique index on balances(user_id, currency);