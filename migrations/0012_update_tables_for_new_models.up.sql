CREATE TABLE IF NOT EXISTS public.care_log (
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    was_watered bool NOT NULL DEFAULT false,
    was_fertilized bool NOT NULL DEFAULT false,
    notes text DEFAULT ''::text,
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    plant_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS public.nectar_users (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    authId text NOT NULL,
    first_name text NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    profile_image text DEFAULT ''::text
);