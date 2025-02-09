-- +goose Up
-- +goose StatementBegin
ALTER TABLE `short_urls` RENAME COLUMN `expired_at` TO `expire_at`;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE `short_urls` RENAME COLUMN `expire_at` TO `expired_at`;
-- +goose StatementEnd
