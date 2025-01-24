-- +goose Up
CREATE TABLE rotations (
    id              bigint primary key GENERATED ALWAYS AS IDENTITY,
    banner_id       varchar(255) references banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    slot_id         varchar(255) references slots(id) ON DELETE CASCADE ON UPDATE CASCADE,
    group_id        varchar(255) references groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
    show_count      bigint not null DEFAULT 1,
    click_count     bigint not null DEFAULT 1
);

CREATE INDEX rotations_on_slot_group_idx ON rotations (slot_id, group_id);
CREATE UNIQUE INDEX rotations_on_banner_slot_group_idx ON rotations (banner_id, slot_id, group_id);

-- +goose Down
DROP INDEX rotations_on_banner_slot_group_idx;
DROP INDEX rotations_on_slot_group_idx;
DROP TABLE rotations;