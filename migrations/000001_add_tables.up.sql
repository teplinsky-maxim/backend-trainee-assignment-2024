BEGIN;

create table if not exists banners
(
    id         bigserial primary key,
    title      text,
    text       text,
    url        text,
    feature_id bigint
);

create table if not exists banner_tags
(
    id        bigserial primary key,
    banner_id bigint,
    tag_id    bigint
);

COMMIT;