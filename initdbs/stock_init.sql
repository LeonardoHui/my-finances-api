CREATE TABLE stock
(
    id bigint NOT NULL,
    name text COLLATE pg_catalog."default",
    CONSTRAINT stock_pkey PRIMARY KEY (id)
);

INSERT INTO stock(id, name) VALUES
 (1, 'D'),
 (2, 'E'),
 (3, 'F');