CREATE TABLE IF NOT EXISTS clicks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id     UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at  TIMESTAMPTZ DEFAULT NOW(),
    referrer    TEXT,
    user_agent  TEXT,
    country     TEXT,
    city        TEXT
);