-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
 		CREATE TABLE IF NOT EXISTS "other_polygon_structure"(
            gid integer NOT NULL DEFAULT nextval('other_polygon_structure_gid_seq'::regclass),
            shape__are numeric,
            shape__len numeric,
            globalid character varying(38) COLLATE pg_catalog."default",
            creationda date,
            creator character varying(128) COLLATE pg_catalog."default",
            editdate date,
            editor character varying(128) COLLATE pg_catalog."default",
            exact_use character varying(254) COLLATE pg_catalog."default",
            usage1 character varying(254) COLLATE pg_catalog."default",
            developmen character varying(254) COLLATE pg_catalog."default",
            structure_ character varying(254) COLLATE pg_catalog."default",
            ghanapostg character varying(254) COLLATE pg_catalog."default",
            street_nam character varying(254) COLLATE pg_catalog."default",
            name character varying(254) COLLATE pg_catalog."default",
            parcel_id character varying(254) COLLATE pg_catalog."default",
            remarks character varying(254) COLLATE pg_catalog."default",
            other_info character varying(254) COLLATE pg_catalog."default",
            other_in_1 character varying(254) COLLATE pg_catalog."default",
            mixed_usag character varying(254) COLLATE pg_catalog."default",
            geom geometry(MultiPolygonZM),
            CONSTRAINT other_polygon_structure_pkey PRIMARY KEY (gid)
        )

            TABLESPACE pg_default;

            ALTER TABLE IF EXISTS public.other_polygon_structure
            OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE other_polygon_structure
-- +goose StatementEnd
