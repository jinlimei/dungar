alter table manic_minute_runtime_state
    add column if not exists cooldown_until timestamp with time zone;
