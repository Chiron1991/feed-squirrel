-- +migrate Up
CREATE TABLE feeds
(
    -- metadata
    id           bigint generated always as identity primary key,
    created_at   timestamptz not null default (now() at time zone 'utc'),
    deleted_at   timestamptz,
    -- values
    last_scraped timestamptz, -- todo: type ok?
    -- see https://github.com/mmcdole/gofeed#default-mappings
    title        text        not null default '',
    description  text        not null default '',
    link         text        not null default '',
    feed_link    text        not null default '',
    updated      timestamptz, -- todo: type ok?
    published    timestamptz, -- todo: type ok?
    author       text        not null default '',
    language     text        not null default '',
    image        text        not null default '',
    copyright    text        not null default '',
    generator    text        not null default '',
    categories   text        not null default ''
);
CREATE TABLE feed_items
(
    -- metadata
    id          bigint generated always as identity primary key,
    created_at  timestamptz not null default (now() at time zone 'utc'),
    deleted_at  timestamptz,
    -- relations
    feed_id     int         not null references feeds (id) on delete cascade,
    -- values, see https://github.com/mmcdole/gofeed#default-mappings
    title       text        not null default '',
    description text        not null default '',
    content     text        not null default '',
    link        text        not null default '',
    updated     timestamptz, -- todo: type ok?
    published   timestamptz, -- todo: type ok?
--     author      text,
    guid        text        not null unique
);

-- +migrate Down
DROP TABLE feed_items;
DROP TABLE feeds;
