ALTER TABLE subjects
ADD COLUMN parent_id UUID;

ALTER TABLE subjects
ADD CONSTRAINT fk_subjects_parent_id FOREIGN KEY (parent_id) REFERENCES subjects (id);
