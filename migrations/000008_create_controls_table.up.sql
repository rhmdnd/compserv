CREATE TABLE IF NOT EXISTS controls (
  id uuid PRIMARY KEY,
  name VARCHAR(255),
  severity VARCHAR(50),
  profile_id UUID,
  metadata_id UUID,
  CONSTRAINT fk_controls_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata (id),
  CONSTRAINT fk_controls_profile_id FOREIGN KEY (profile_id) REFERENCES profiles (id)
);
