-- +goose Up
-- +goose StatementBegin
CREATE UNLOGGED TABLE IF NOT EXISTS tbl_cache
(
    key character varying(100) NOT NULL,
    value text,
    created timestamp with time zone,
    PRIMARY KEY (key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "tbl_cache";
-- +goose StatementEnd
