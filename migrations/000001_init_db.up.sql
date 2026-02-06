CREATE TABLE users 
(
	id varchar(36) unique,
    display_name varchar(255) not null,
    username varchar(255) not null,
    email varchar(255) not null,
    password varchar(2048) not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT uni_username UNIQUE (username),
    CONSTRAINT uni_email UNIQUE (email)
);
