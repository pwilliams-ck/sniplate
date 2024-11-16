CREATE INDEX IF NOT EXISTS snips_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS snips_tags_idx ON movies USING GIN (tags);
