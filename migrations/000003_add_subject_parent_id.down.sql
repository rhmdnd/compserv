ALTER TABLE subjects DROP CONSTRAINT fk_subjects_parent_id;

ALTER TABLE subjects DROP COLUMN parent_id;
