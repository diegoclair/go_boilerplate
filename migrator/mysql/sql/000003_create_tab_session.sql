CREATE TABLE IF NOT EXISTS tab_session (
    session_id INT NOT NULL AUTO_INCREMENT,
    session_uuid CHAR(36) NOT NULL,
    account_id INT NULL,
    refresh_token VARCHAR(1500) NOT NULL,
    user_agent VARCHAR(1000) NOT NULL,
    client_ip VARCHAR(500) NOT NULL,
    is_blocked boolean NOT NULL DEFAULT false,
    refresh_token_expires_at TIMESTAMP NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (session_id),
    UNIQUE INDEX session_id_UNIQUE (session_id ASC) VISIBLE,
    INDEX fk_tab_session_tab_account_idx (account_id ASC) VISIBLE,

    CONSTRAINT fk_tab_session_tab_account
        FOREIGN KEY (account_id)
        REFERENCES tab_account (account_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET=utf8;