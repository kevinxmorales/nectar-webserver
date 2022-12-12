CREATE TABLE IF NOT EXISTS public.plant (
     id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
     user_id uuid NOT NULL,
     common_name text NOT NULL,
     scientific_name text,
     toxicity text,
     created_at timestamp with time zone default current_timestamp
);

CREATE TABLE IF NOT EXISTS public.plant_search_terms (
    plant_id uuid NOT NULL,
    term text NOT NULL
);

