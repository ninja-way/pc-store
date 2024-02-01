CREATE TABLE refresh_tokens (
    id serial primary key NOT NULL,
    user_id integer NOT NULL,
    token character varying(255) NOT NULL,
    expires_at timestamp without time zone NOT NULL
);
