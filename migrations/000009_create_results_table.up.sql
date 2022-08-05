CREATE TABLE IF NOT EXISTS results (
        id uuid PRIMARY KEY,
        name VARCHAR(255),
        outcome VARCHAR(255),
        instruction TEXT,
        rationale TEXT,
        control_id UUID,
        metadata_id UUID,
        subject_id UUID,
        assessment_id UUID,
        CONSTRAINT fk_results_control_id FOREIGN KEY (control_id) REFERENCES controls (id),
        CONSTRAINT fk_results_metadata_id FOREIGN KEY (metadata_id) REFERENCES metadata (id),
        CONSTRAINT fk_results_subject_id FOREIGN KEY (subject_id) REFERENCES subjects (id),
        CONSTRAINT fk_results_assessment_id FOREIGN KEY (assessment_id) REFERENCES assessments (id)
);
