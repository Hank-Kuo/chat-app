CREATE TABLE IF NOT EXISTS reply (
    channel_id UUID,
    message_id BIGINT,
    reply_id BIGINT,
    content TEXT,
    user_id UUID,
    username TEXT,
    created_at  TIMESTAMP,
    PRIMARY KEY(message_id, reply_id)
) WITH CLUSTERING ORDER BY (reply_id DESC);