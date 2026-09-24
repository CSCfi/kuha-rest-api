--
-- Archinisis schema (from empty_archinisis_dump.sql, 2026-09-24)
-- Source of truth for sqlc. Data COPY blocks and sequence setval calls removed.
--

CREATE TABLE public.athlete (
    national_id character varying(50) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    initials character varying(10),
    date_of_birth date,
    height numeric,
    weight numeric
);

CREATE TABLE public.measurement_group (
    measurement_group_id integer NOT NULL
);

CREATE SEQUENCE public.measurement_group_measurement_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.measurement_group_measurement_group_id_seq OWNED BY public.measurement_group.measurement_group_id;

CREATE TABLE public.measurement (
    measurement_group_id integer,
    measurement_id integer NOT NULL,
    national_id character varying(50),
    discipline character varying(100),
    session_name character varying(255),
    place character varying(100),
    race_id integer,
    start_time timestamp without time zone,
    stop_time timestamp without time zone,
    nb_segments integer,
    comment text
);

CREATE SEQUENCE public.measurement_measurement_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.measurement_measurement_id_seq OWNED BY public.measurement.measurement_id;

CREATE TABLE public.report (
    report_id integer NOT NULL,
    sportti_id character varying(50),
    session_id integer,
    race_report text
);

CREATE SEQUENCE public.report_report_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.report_report_id_seq OWNED BY public.report.report_id;

ALTER TABLE ONLY public.measurement ALTER COLUMN measurement_id SET DEFAULT nextval('public.measurement_measurement_id_seq'::regclass);
ALTER TABLE ONLY public.measurement_group ALTER COLUMN measurement_group_id SET DEFAULT nextval('public.measurement_group_measurement_group_id_seq'::regclass);
ALTER TABLE ONLY public.report ALTER COLUMN report_id SET DEFAULT nextval('public.report_report_id_seq'::regclass);

ALTER TABLE ONLY public.athlete
    ADD CONSTRAINT athlete_pkey PRIMARY KEY (national_id);

ALTER TABLE ONLY public.measurement_group
    ADD CONSTRAINT measurement_group_pkey PRIMARY KEY (measurement_group_id);

ALTER TABLE ONLY public.measurement
    ADD CONSTRAINT measurement_pkey PRIMARY KEY (measurement_id);

ALTER TABLE ONLY public.report
    ADD CONSTRAINT report_pkey PRIMARY KEY (report_id);

ALTER TABLE ONLY public.measurement
    ADD CONSTRAINT measurement_measurement_group_id_fkey FOREIGN KEY (measurement_group_id) REFERENCES public.measurement_group(measurement_group_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.measurement
    ADD CONSTRAINT measurement_national_id_fkey FOREIGN KEY (national_id) REFERENCES public.athlete(national_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.report
    ADD CONSTRAINT report_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.measurement_group(measurement_group_id) ON DELETE CASCADE;

ALTER TABLE ONLY public.report
    ADD CONSTRAINT report_sportti_id_fkey FOREIGN KEY (sportti_id) REFERENCES public.athlete(national_id) ON DELETE CASCADE;
