CREATE TABLE IF NOT EXISTS public.user_credentials (
   id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
   user_id uuid NOT NULL,
   password_digest text NOT NULL,
   last_updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);