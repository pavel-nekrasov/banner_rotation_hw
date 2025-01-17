-- +goose Up
CREATE table banners (
    id              varchar(255) primary key,
    description     text
);

CREATE table slots (
    id              varchar(255) primary key,
    description     text
);

CREATE table groups (
    id              varchar(255) primary key,
    description     text
);

-- +goose Down
drop table banners;
drop table slots;
drop table groups;