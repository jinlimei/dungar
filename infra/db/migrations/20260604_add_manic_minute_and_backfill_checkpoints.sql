create table if not exists manic_minute_events
(
    id serial not null
        constraint manic_minute_events_pk
            primary key,
    server_id varchar(32) not null,
    channel_id varchar(32) not null,
    trigger_word varchar(128) not null,
    trigger_message_id varchar(32),
    triggered_by_user_id varchar(32),
    status varchar(20) not null,
    stop_reason varchar(50),
    message_count integer default 0 not null,
    started_at timestamp with time zone not null,
    ended_at timestamp with time zone
);

create index if not exists manic_minute_events_server_started_idx
    on manic_minute_events (server_id, started_at desc);

create index if not exists manic_minute_events_channel_started_idx
    on manic_minute_events (server_id, channel_id, started_at desc);

create index if not exists manic_minute_events_word_started_idx
    on manic_minute_events (server_id, trigger_word, started_at desc);

create table if not exists discord_backfill_checkpoints
(
    channel_id varchar(32) not null
        constraint discord_backfill_checkpoints_pk
            primary key,
    guild_id varchar(32) not null,
    channel_name varchar(100),
    status varchar(20) not null,
    before_message_id varchar(32),
    since_ts timestamp with time zone not null,
    fetched integer default 0 not null,
    stored integer default 0 not null,
    skipped_duplicate integer default 0 not null,
    skipped_unusable integer default 0 not null,
    skipped_old integer default 0 not null,
    last_message_timestamp timestamp with time zone,
    last_error text,
    started_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    completed_at timestamp with time zone
);

create index if not exists discord_backfill_checkpoints_status_idx
    on discord_backfill_checkpoints (status, updated_at desc);
