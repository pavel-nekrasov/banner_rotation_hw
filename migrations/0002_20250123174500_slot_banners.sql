-- +goose Up
CREATE TABLE slot_banners (
    id              bigint primary key GENERATED ALWAYS AS IDENTITY,
    banner_id       varchar(255) references banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    slot_id         varchar(255) references slots(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX slot_banners_on_slot_idx ON slot_banners (slot_id);
CREATE UNIQUE INDEX slot_banners_on_banner_slot_idx ON slot_banners (banner_id, slot_id);

-- +goose Down
DROP INDEX slot_banners_on_banner_slot_idx;
DROP INDEX slot_banners_on_slot_idx;
DROP TABLE slot_banners;
