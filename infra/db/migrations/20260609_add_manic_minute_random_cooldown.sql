alter table manic_minute_events
    add column if not exists cooldown_until timestamp with time zone;

update manic_minute_events
set cooldown_until = started_at + interval '120 minutes'
where cooldown_until is null;

alter table manic_minute_events
    alter column cooldown_until set not null;

create index if not exists manic_minute_events_channel_cooldown_idx
    on manic_minute_events (server_id, channel_id, cooldown_until desc);

create index if not exists manic_minute_events_word_cooldown_idx
    on manic_minute_events (server_id, trigger_word, cooldown_until desc);
