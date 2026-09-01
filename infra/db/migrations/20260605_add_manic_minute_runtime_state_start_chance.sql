alter table manic_minute_runtime_state
    add column if not exists start_chance double precision default 0.20 not null;
