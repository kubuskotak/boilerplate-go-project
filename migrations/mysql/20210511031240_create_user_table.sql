-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users`
(
    `username`            varchar(140)       NOT NULL,
    `hashed_password`     varchar(140)       NOT NULL,
    `full_name`           varchar(140)       NOT NULL,
    `email`               varchar(40) UNIQUE NOT NULL,
    `password_changed_at` Timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at`          Timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`username`)
)
    CHARACTER SET = utf8
    COLLATE = utf8_general_ci
    ENGINE = InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
