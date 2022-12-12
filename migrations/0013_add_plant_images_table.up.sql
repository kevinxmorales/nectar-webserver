CREATE TABLE IF NOT EXISTS public.plant_images (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    image text NOT NULL,
    plant_id uuid NOT NULL
);