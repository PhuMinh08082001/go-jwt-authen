create table if not exists users (
                                     id bigserial primary key,
                                     user_name varchar(128),
                                    password varchar(128)
    );
