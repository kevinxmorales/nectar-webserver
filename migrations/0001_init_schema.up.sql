CREATE TABLE IF NOT EXISTS public.users (
     id serial PRIMARY KEY ,
     usr_frst_nm text,
     usr_lst_nm text,
     usr_email text UNIQUE NOT NULL,
     usr_psswrd text,
     usr_username text
);

CREATE TABLE IF NOT EXISTS public.plants (
    id serial PRIMARY KEY,
    plnt_nm text,
    plnt_usr_id int4,
    plnt_ctgry_id integer,
    plnt_urls jsonb DEFAULT '{}',
    plnt_created_at timestamp with time zone default current_timestamp
);

CREATE TABLE IF NOT EXISTS public.plant_category (
    ctgry_id serial PRIMARY KEY,
    ctgry_clr text,
    ctgry_icon text,
    ctgry_lbl text
);

CREATE TABLE "public"."plant_care_log" (
    "pcl_date" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "pcl_was_watered" bool NOT NULL DEFAULT false,
    "pcl_was_fertilized" bool NOT NULL DEFAULT false,
    "pcl_notes" text DEFAULT ''::text,
    "id" serial NOT NULL,
    "pcl_plnt_id" int4 NOT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO public.users (usr_frst_nm, usr_lst_nm, usr_psswrd, usr_email) VALUES ('Kevin', 'Morales','$2a$04$KEpBH4ODt4u2oZB8M8J.7eZ8ucj91u6HGQ3bw89NyNuaRr9b32Zeu', 'kevin@domain.com') ON CONFLICT DO NOTHING;