-- +goose Up
CREATE TABLE IF NOT EXISTS tab_account (
    account_id SERIAL PRIMARY KEY,
    account_uuid UUID NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    name VARCHAR(450) NOT NULL,
    secret VARCHAR(200) NOT NULL,
    balance DECIMAL(7,2) DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT cpf_unique UNIQUE (cpf)
);

CREATE TABLE IF NOT EXISTS tab_transfer (
    transfer_id SERIAL PRIMARY KEY,
    account_origin_id INT NOT NULL,
    account_destination_id INT NOT NULL,
    amount DECIMAL(7,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tab_transfer_tab_account
        FOREIGN KEY (account_origin_id)
        REFERENCES tab_account (account_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT fk_tab_transfer_tab_account1
        FOREIGN KEY (account_destination_id)
        REFERENCES tab_account (account_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);

CREATE INDEX idx_tab_transfer_account_origin ON tab_transfer (account_origin_id);
CREATE INDEX idx_tab_transfer_account_destination ON tab_transfer (account_destination_id);

-- +goose Down
DROP TABLE IF EXISTS tab_transfer;
DROP TABLE IF EXISTS tab_account;
