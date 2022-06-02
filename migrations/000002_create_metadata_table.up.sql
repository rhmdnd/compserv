CREATE TABLE IF NOT EXISTS metadata (
   id uuid PRIMARY KEY,
   created_at timestamp without time zone,
   updated_at timestamp without time zone,
   version VARCHAR(50),
   description TEXT
);
