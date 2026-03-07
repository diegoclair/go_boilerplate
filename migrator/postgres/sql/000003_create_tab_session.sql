-- +goose Up
CREATE TABLE IF NOT EXISTS tab_session (
    session_id SERIAL PRIMARY KEY,
    session_uuid UUID NOT NULL,
    account_id INT,
    refresh_token VARCHAR(1500) NOT NULL,
    user_agent VARCHAR(1000) NOT NULL,
    client_ip VARCHAR(500) NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    update_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tab_session_tab_account
        FOREIGN KEY (account_id)
        REFERENCES tab_account (account_id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);

CREATE INDEX idx_tab_session_account ON tab_session (account_id);

-- +goose Down
DROP TABLE IF EXISTS tab_session;
