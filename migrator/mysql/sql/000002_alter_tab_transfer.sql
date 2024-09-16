ALTER TABLE tab_transfer
    ADD COLUMN transfer_uuid CHAR(36) NOT NULL after transfer_id;