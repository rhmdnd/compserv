CREATE TABLE IF NOT EXISTS profiles (
  id uuid PRIMARY KEY,
  name VARCHAR(255),
  metadata_id UUID,
  catalog_id UUID,
  CONSTRAINT fk_profiles_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata(id),
  CONSTRAINT fk_profiles_catalog_id FOREIGN KEY (catalog_id) REFERENCES catalogs(id)
);

