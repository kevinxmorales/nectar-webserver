CREATE TABLE IF NOT EXISTS plants (
    plnt_id uuid PRIMARY KEY,
    plnt_nm text,
    plnt_usr_id uuid,
    plnt_ctgry_id integer,
    plnt_imgs jsonb,
    plnt_created_at timestamp with time zone default current_timestamp
);

CREATE TABLE IF NOT EXISTS plant_category (
    ctgry_id serial PRIMARY KEY,
    ctgry_clr text,
    ctgry_icon text,
    ctgry_lbl text
);

CREATE TABLE IF NOT EXISTS users (
    usr_id uuid,
    usr_frst_nm text,
    usr_lst_nm text,
    usr_email text UNIQUE NOT NULL,
    usr_psswrd text
);

INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('df0cfbfb-5bd1-44e8-a526-7b3121be276d', 'Monstera', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 5, '{"files": [{"fileName": "monstera"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('6a5112be-6c83-43eb-925b-f20a8ac848e5', 'Philodendron', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 1, '{"files": [{"fileName": "plant"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('db139401-5c91-441f-9887-e766cbdac1b8', 'Poison1', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 1, '{"files": [{"fileName": "poison"}, {"fileName": "poison"}, {"fileName": "poison"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('67e353c6-6d49-4978-9a44-9dbe81904b97', 'Ficus', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 5, '{"files": [{"fileName": "ficus"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('414cb52d-66ed-47c0-8400-e0092c5f7f79', 'Another Monstera', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 3, '{"files": [{"fileName": "monstera"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('e8b3657b-5d66-445d-81d2-1593bad5e996', 'Another Plant', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 3, '{"files": [{"fileName": "plant"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('6b299c70-df08-482a-bbd6-4f45a9956be6', 'Philo', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 1, '{"files": [{"fileName": "plant"}]}');
INSERT INTO public.plants (plnt_id, plnt_nm, plnt_usr_id, plnt_ctgry_id, plnt_imgs) VALUES ('25f51176-377f-40cc-a530-2e60488beca0', 'Ficus 2', '4d609f18-de3d-485a-b8b8-23e95b9c76f8', 5, '{"files": [{"fileName": "ficus"}]}');

INSERT INTO public.users (usr_id, usr_frst_nm, usr_lst_nm, usr_psswrd, usr_email) VALUES ('4d609f18-de3d-485a-b8b8-23e95b9c76f8', 'Kevin', 'Morales','$2a$04$KEpBH4ODt4u2oZB8M8J.7eZ8ucj91u6HGQ3bw89NyNuaRr9b32Zeu', 'kevin@domain.com');