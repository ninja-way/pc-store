CREATE TABLE logs (
    id serial primary key NOT NULL,
    entity character varying(255),
    action character varying(255),
    entity_id integer,
    timestamp timestamp without time zone DEFAULT now() NOT NULL
);