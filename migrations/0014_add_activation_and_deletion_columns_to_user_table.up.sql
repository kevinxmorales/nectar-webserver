ALTER TABLE nectar_users
ADD COLUMN account_creation_date timestamp DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN account_deletion_date timestamp DEFAULT 'infinity'