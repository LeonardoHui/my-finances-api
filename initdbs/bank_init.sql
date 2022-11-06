CREATE TABLE bank
(
    id bigint NOT NULL,
    name text COLLATE pg_catalog."default",
    CONSTRAINT bank_pkey PRIMARY KEY (id)
);

INSERT INTO bank(id, name) VALUES
 (1, 'A'),
 (2, 'B'),
 (3, 'C');