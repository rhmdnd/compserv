CREATE TABLE IF NOT EXISTS catalogs (
  id uuid PRIMARY KEY,
  name VARCHAR(255),
  metadata_id UUID,
  content TEXT,
  CONSTRAINT fk_catalogs_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata (id)
);
