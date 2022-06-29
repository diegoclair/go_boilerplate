package migrations

import "github.com/GuiaBolso/darwin"

//Only work doing 1 command per version, you cannot create two tables in the same script, need to create a new version
var (
	Migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "Create tab_account",
			Script: `
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
				UNIQUE INDEX cpf_UNIQUE (cpf ASC) VISIBLE)
				ENGINE = InnoDB CHARACTER SET=utf8;
			`,
		},
		{
			Version:     2,
			Description: "Create tab_transfer",
			Script: `
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
					ON UPDATE NO ACTION)
				ENGINE = InnoDB CHARACTER SET=utf8;
			`,
		},
		{
			Version:     3,
			Description: "Add uuid field into tab_transfer",
			Script: `
				ALTER TABLE tab_transfer
					ADD COLUMN transfer_uuid CHAR(36) NOT NULL after transfer_id;
			`,
		},
		{
			Version:     4,
			Description: "Create tab session",
			Script: `
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
				
				CONSTRAINT fk_tab_session_tab_user
					FOREIGN KEY (account_id)
					REFERENCES tab_account (account_id)
					ON DELETE NO ACTION
					ON UPDATE NO ACTION)
				ENGINE = InnoDB CHARACTER SET=utf8;
			`,
		},
	}
)
