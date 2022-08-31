BEGIN;
CREATE TABLE IF NOT EXISTS following (
    f_user_id int4,
    f_user_being_followed_id int4,
    f_started_following_date timestamp with time zone default current_timestamp
);
COMMIT;