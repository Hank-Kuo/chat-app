CREATE TABLE IF NOT EXISTS message (
    channel_id UUID,
    bucket INT,
    message_id BIGINT,
    content TEXT,
    user_id UUID,
    username TEXT,
    created_at  TIMESTAMP,
    PRIMARY KEY((channel_id, bucket), message_id)
)  WITH CLUSTERING ORDER BY (message_id DESC);