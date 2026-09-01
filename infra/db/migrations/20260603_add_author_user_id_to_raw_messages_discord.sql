ALTER TABLE raw_messages_discord
    ADD COLUMN IF NOT EXISTS author_user_id varchar(32);

CREATE INDEX IF NOT EXISTS raw_messages_discord_author_user_idx
    ON raw_messages_discord (author_user_id);
