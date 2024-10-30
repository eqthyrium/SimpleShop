

-- DROP table if exists user;
-- DROP table if exists post;
-- DROP table if exists comment;
-- DROP table if exists category;
-- Drop table if exists post-category;


CREATE table user (
    user_id integer primary key,
    nickname varchar(50),
    email varchar(50),
    password varchar(256),
    role varchar(15)
);
