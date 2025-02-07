-- +goose Up
-- +goose StatementBegin
CREATE TABLE `short_urls` (
	`id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
	`url` TEXT NOT NULL,
	`target_id` CHAR(8) NOT NULL,
	`expired_at` DATETIME NOT NULL,
	`created_at` DATETIME NOT NULL,

  	UNIQUE INDEX `uqidx_target_id` (`target_id`)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `short_urls`;
-- +goose StatementEnd
