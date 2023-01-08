ALTER TABLE plant
    ADD COLUMN creation_date timestamp DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN last_update_date timestamp DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN deletion_date timestamp DEFAULT 'infinity'