CREATE TABLE users (
    id serial primary key NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    registered_at timestamp without time zone DEFAULT now() NOT NULL
);