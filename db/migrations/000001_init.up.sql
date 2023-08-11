CREATE TABLE users (
    id int primary key not NULL,
    banned boolean default false not null,
    followed boolean default true
);

CREATE TABLE track (
    id serial,
    user_id int REFERENCES users(id) on delete cascade,
    stock varchar not null,
    market varchar not null,
    board_group int not null,
    security varchar not null,
    tracked_volume int not null,
    isTracked boolean default true
);