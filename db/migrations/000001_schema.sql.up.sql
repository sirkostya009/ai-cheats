drop table customers;

create table customers (
    id               int primary key,
    telegram         text unique,
    active           bool,
    ips              inet[],
    max_ips          int,
    model            text,
    requests         int,
    request_tokens   int,
    generated_tokens int
);
