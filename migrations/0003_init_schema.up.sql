BEGIN;
ALTER TABLE IF EXISTS users ADD COLUMN  "usr_image_url" text;
COMMIT;
