-- +goose Up
ALTER TABLE tab_transfer
    ADD COLUMN transfer_uuid UUID NOT NULL;

-- +goose Down
ALTER TABLE tab_transfer
    DROP COLUMN transfer_uuid;
