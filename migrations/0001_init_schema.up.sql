CREATE TABLE IF NOT EXISTS plants (
    plnt_id uuid PRIMARY KEY,
    plnt_nm text,
    plnt_usr_id uuid,
    plnt_ctgry_id integer
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

