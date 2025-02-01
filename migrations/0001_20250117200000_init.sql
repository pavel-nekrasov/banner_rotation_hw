-- +goose Up
CREATE TABLE banners (
    id              varchar(255) primary key,
    description     text
);

CREATE TABLE groups (
    id              varchar(255) primary key,
    description     text
);

CREATE TABLE slots (
    id              varchar(255) primary key,
    description     text
);

-- +goose Down
DROP TABLE slots;
DROP TABLE groups;
DROP TABLE banners;