CREATE TABLE IF NOT EXISTS raw_message_user_blocklist
(
    protocol_driver varchar(20) NOT NULL,
    server_id varchar(32) NOT NULL,
    user_id varchar(32) NOT NULL,
    nick varchar(100) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT raw_message_user_blocklist_pk
        PRIMARY KEY (protocol_driver, server_id, user_id)
);

CREATE INDEX IF NOT EXISTS raw_message_user_blocklist_lookup_idx
    ON raw_message_user_blocklist (protocol_driver, server_id, user_id);
