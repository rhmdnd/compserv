CREATE TABLE IF NOT EXISTS assessments (
  id uuid PRIMARY KEY,
  name VARCHAR(255),
  metadata_id UUID,
  CONSTRAINT fk_assessments_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata (id)
);
