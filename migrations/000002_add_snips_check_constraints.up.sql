ALTER TABLE snips ADD CONSTRAINT tags_length_check CHECK (array_length(tags, 1) BETWEEN 0 AND 10);

