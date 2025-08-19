-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tbl_key_value" (
	"key"	TEXT,
	"value"	TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tbl_key_value";
-- +goose StatementEnd
