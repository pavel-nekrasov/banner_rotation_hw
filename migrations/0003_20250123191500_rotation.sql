-- +goose Up
CREATE TABLE banner_stats (
    id              bigint primary key GENERATED ALWAYS AS IDENTITY,
    banner_id       varchar(255) references banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    slot_id         varchar(255) references slots(id) ON DELETE CASCADE ON UPDATE CASCADE,
    group_id        varchar(255) references groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
    show_count      bigint not null DEFAULT 1,
    click_count     bigint not null DEFAULT 1
);

CREATE INDEX banner_stats_on_slot_group_idx ON banner_stats (slot_id, group_id);
CREATE UNIQUE INDEX banner_stats_on_banner_slot_group_idx ON banner_stats (banner_id, slot_id, group_id);

-- +goose Down
DROP INDEX banner_stats_on_banner_slot_group_idx;
DROP INDEX banner_stats_on_slot_group_idx;
DROP TABLE banner_stats;