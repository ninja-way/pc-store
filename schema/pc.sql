CREATE TABLE pc (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    cpu character varying(255) NOT NULL,
    videocard character varying(255),
    ram integer NOT NULL,
    data_storage character varying(255),
    added_at timestamp without time zone DEFAULT now() NOT NULL,
    price integer NOT NULL
);