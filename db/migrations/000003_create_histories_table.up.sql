create table if not exists histories(
    id bigserial primary key,
    user_id bigint not null references users(id),
    balance numeric(20, 2) not null,
    currency varchar(3) not null,
    transfer_proof_image varchar not null,
    source  jsonb not null default '{}'::jsonb,
    created_at bigint not null default extract(epoch from current_timestamp)
);

-- add index to histories table
create index on histories(user_id);



