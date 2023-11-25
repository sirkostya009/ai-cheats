create table if not exists customers
(
    id       serial primary key,
    telegram text unique,
    active   bool   default true,
    hashes   text[] default '{}'::text[],
    max_ips  int    default 1,
    model    text   default 'gpt-3.5-turbo'
);

create table if not exists requests
(
    customer_id       int not null references customers (id),
    created_at        timestamp,
    finished_at       timestamp default now(),
    completion_tokens int,
    prompt_tokens     int,
    status            int,
    reason            text,
    model             text
);
