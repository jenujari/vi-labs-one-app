-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS tbl_seven_fifty
ADD COLUMN instrument_token integer;
ALTER TABLE IF EXISTS tbl_seven_fifty
ADD COLUMN exchange_token integer;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE tbl_seven_fifty DROP COLUMN instrument_token;
ALTER TABLE tbl_seven_fifty DROP COLUMN exchange_token;
-- +goose StatementEnd