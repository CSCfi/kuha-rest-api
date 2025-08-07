--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg110+2)
-- Dumped by pg_dump version 17.3 (Ubuntu 17.3-3.pgdg24.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: athlete; Type: TABLE; Schema: public; Owner: archadmin
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


ALTER TABLE public.athlete OWNER TO archadmin;

--
-- Name: measurement; Type: TABLE; Schema: public; Owner: archadmin
--

CREATE TABLE public.measurement (
    measurement_group_id integer NOT NULL,
    measurement_id integer,
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


ALTER TABLE public.measurement OWNER TO archadmin;

--
-- Name: mittaus_measurement_group_id_seq; Type: SEQUENCE; Schema: public; Owner: archadmin
--

CREATE SEQUENCE public.mittaus_measurement_group_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mittaus_measurement_group_id_seq OWNER TO archadmin;

--
-- Name: mittaus_measurement_group_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: archadmin
--

ALTER SEQUENCE public.mittaus_measurement_group_id_seq OWNED BY public.measurement.measurement_group_id;


--
-- Name: report; Type: TABLE; Schema: public; Owner: archadmin
--

CREATE TABLE public.report (
    report_id integer NOT NULL,
    sportti_id character varying(50),
    session_id integer,
    race_report text
);


ALTER TABLE public.report OWNER TO archadmin;

--
-- Name: raportti_report_id_seq; Type: SEQUENCE; Schema: public; Owner: archadmin
--

CREATE SEQUENCE public.raportti_report_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.raportti_report_id_seq OWNER TO archadmin;

--
-- Name: raportti_report_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: archadmin
--

ALTER SEQUENCE public.raportti_report_id_seq OWNED BY public.report.report_id;


--
-- Name: sportti_id_list; Type: TABLE; Schema: public; Owner: archadmin
--

CREATE TABLE public.sportti_id_list (
    sportti_id character varying(64) NOT NULL
);


ALTER TABLE public.sportti_id_list OWNER TO archadmin;

--
-- Name: measurement measurement_group_id; Type: DEFAULT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.measurement ALTER COLUMN measurement_group_id SET DEFAULT nextval('public.mittaus_measurement_group_id_seq'::regclass);


--
-- Name: report report_id; Type: DEFAULT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.report ALTER COLUMN report_id SET DEFAULT nextval('public.raportti_report_id_seq'::regclass);

--
-- Name: mittaus_measurement_group_id_seq; Type: SEQUENCE SET; Schema: public; Owner: archadmin
--

SELECT pg_catalog.setval('public.mittaus_measurement_group_id_seq', 1, false);


--
-- Name: raportti_report_id_seq; Type: SEQUENCE SET; Schema: public; Owner: archadmin
--

SELECT pg_catalog.setval('public.raportti_report_id_seq', 1, false);


--
-- Name: measurement mittaus_pkey; Type: CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.measurement
    ADD CONSTRAINT mittaus_pkey PRIMARY KEY (measurement_group_id);


--
-- Name: report raportti_pkey; Type: CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.report
    ADD CONSTRAINT raportti_pkey PRIMARY KEY (report_id);


--
-- Name: report raportti_session_id_key; Type: CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.report
    ADD CONSTRAINT raportti_session_id_key UNIQUE (session_id);


--
-- Name: sportti_id_list sportti_id_list_pkey; Type: CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.sportti_id_list
    ADD CONSTRAINT sportti_id_list_pkey PRIMARY KEY (sportti_id);


--
-- Name: athlete urheilija_pkey; Type: CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.athlete
    ADD CONSTRAINT urheilija_pkey PRIMARY KEY (national_id);


--
-- Name: measurement mittaus_national_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.measurement
    ADD CONSTRAINT mittaus_national_id_fkey FOREIGN KEY (national_id) REFERENCES public.athlete(national_id) ON DELETE CASCADE;


--
-- Name: report raportti_session_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.report
    ADD CONSTRAINT raportti_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.measurement(measurement_group_id) ON DELETE CASCADE;


--
-- Name: report raportti_sportti_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: archadmin
--

ALTER TABLE ONLY public.report
    ADD CONSTRAINT raportti_sportti_id_fkey FOREIGN KEY (sportti_id) REFERENCES public.athlete(national_id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO archadmin;


--
-- PostgreSQL database dump complete
--

