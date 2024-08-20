CREATE TABLE IF NOT EXISTS tab_account (
    account_id INT NOT NULL AUTO_INCREMENT,
    account_uuid CHAR(36) NOT NULL,
    cpf VARCHAR(11) NOT NULL,
    name VARCHAR(450) NOT NULL,
    secret VARCHAR(200) NOT NULL,
    balance DECIMAL(7,2) NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    active TINYINT(1) NOT NULL DEFAULT 1,

    PRIMARY KEY (account_id),
    UNIQUE INDEX account_id_UNIQUE (account_id ASC) VISIBLE,
    UNIQUE INDEX cpf_UNIQUE (cpf ASC) VISIBLE
) ENGINE = InnoDB CHARACTER SET=utf8;

CREATE TABLE IF NOT EXISTS tab_transfer (
    transfer_id INT NOT NULL AUTO_INCREMENT,
    account_origin_id INT NOT NULL,
    account_destination_id INT NOT NULL,
    amount DECIMAL(7,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

PRIMARY KEY (transfer_id),
UNIQUE INDEX transfer_id_UNIQUE (transfer_id ASC) VISIBLE,
INDEX fk_tab_transfer_tab_account_idx (account_origin_id ASC) VISIBLE,
INDEX fk_tab_transfer_tab_account1_idx (account_destination_id ASC) VISIBLE,

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
) ENGINE = InnoDB CHARACTER SET=utf8;