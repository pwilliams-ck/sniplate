CREATE TABLE IF NOT EXISTS snips (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    content text NOT NULL,
    tags text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
