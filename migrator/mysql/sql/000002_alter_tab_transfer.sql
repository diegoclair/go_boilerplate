-- +goose Up
ALTER TABLE tab_transfer
    ADD COLUMN transfer_uuid CHAR(36) NOT NULL after transfer_id;

-- +goose Down
ALTER TABLE tab_transfer
    DROP COLUMN transfer_uuid;
