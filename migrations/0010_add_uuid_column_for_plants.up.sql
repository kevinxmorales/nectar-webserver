CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE plants ADD COLUMN plnt_id uuid NOT NULL DEFAULT uuid_generate_v4();

