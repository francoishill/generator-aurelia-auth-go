-- +migrate Up

# These files contains migration sql scripts for up and down, refer to https://github.com/rubenv/sql-migrate

CREATE TABLE user (
    id              INTEGER         NOT NULL AUTO_INCREMENT,
    activated_on    DATETIME        DEFAULT NULL,
    is_admin        BOOLEAN         DEFAULT 0,
    full_name       VARCHAR(240)    DEFAULT 'NO_FULL_NAME',
    email           VARCHAR(240)    NOT NULL,
    hashed_password VARCHAR(240)    NOT NULL,
    salt            VARCHAR(240)    NOT NULL,
    banned_on       DATETIME        DEFAULT NULL,
    created         DATETIME        DEFAULT CURRENT_TIMESTAMP,
    updated         TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (email)
) ENGINE=InnoDB;



-- +migrate Down

DROP TABLE user;