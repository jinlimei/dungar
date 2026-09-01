create table if not exists manic_minute_runtime_state
(
    protocol_driver varchar(20) not null
        constraint manic_minute_runtime_state_pk
            primary key,
    trigger_word varchar(128) not null,
    active boolean default false not null,
    active_server_id varchar(32),
    active_channel_id varchar(32),
    updated_reason varchar(50),
    updated_at timestamp with time zone not null
);
