ALTER TABLE subjects
Add COLUMN metadata_id UUID;

ALTER TABLE subjects
ADD CONSTRAINT fk_subjects_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata (id);
