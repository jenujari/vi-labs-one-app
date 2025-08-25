CREATE UNLOGGED TABLE IF NOT EXISTS tbl_cache
(
    key character varying(100) NOT NULL,
    value text,
    created timestamp with time zone,
    PRIMARY KEY (key)
);


CREATE TABLE IF NOT EXISTS tbl_seven_fifty
(
    id smallserial NOT NULL,
    symbol character varying(500) NOT NULL,
    full_name text,
    CONSTRAINT prime_id PRIMARY KEY (id),
    CONSTRAINT unique_symbol UNIQUE (symbol)
);