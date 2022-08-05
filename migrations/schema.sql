--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4 (Debian 14.4-1.pgdg110+1)
-- Dumped by pg_dump version 14.4 (Debian 14.4-1.pgdg110+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: assessments; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.assessments (
    id uuid NOT NULL,
    name character varying(255),
    metadata_id uuid
);


ALTER TABLE public.assessments OWNER TO dbadmin;

--
-- Name: catalogs; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.catalogs (
    id uuid NOT NULL,
    name character varying(255),
    metadata_id uuid,
    content text
);


ALTER TABLE public.catalogs OWNER TO dbadmin;

--
-- Name: controls; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.controls (
    id uuid NOT NULL,
    name character varying(255),
    severity character varying(50),
    profile_id uuid,
    metadata_id uuid
);


ALTER TABLE public.controls OWNER TO dbadmin;

--
-- Name: metadata; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.metadata (
    id uuid NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    version character varying(50),
    description text
);


ALTER TABLE public.metadata OWNER TO dbadmin;

--
-- Name: profiles; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.profiles (
    id uuid NOT NULL,
    name character varying(255),
    metadata_id uuid,
    catalog_id uuid
);


ALTER TABLE public.profiles OWNER TO dbadmin;

--
-- Name: results; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.results (
    id uuid NOT NULL,
    name character varying(255),
    outcome character varying(255),
    instruction text,
    rationale text,
    control_id uuid,
    metadata_id uuid,
    subject_id uuid,
    assessment_id uuid
);


ALTER TABLE public.results OWNER TO dbadmin;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO dbadmin;

--
-- Name: subjects; Type: TABLE; Schema: public; Owner: dbadmin
--

CREATE TABLE public.subjects (
    id uuid NOT NULL,
    name character varying(255),
    type character varying(50),
    parent_id uuid,
    metadata_id uuid
);


ALTER TABLE public.subjects OWNER TO dbadmin;

--
-- Name: assessments assessments_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.assessments
    ADD CONSTRAINT assessments_pkey PRIMARY KEY (id);


--
-- Name: catalogs catalogs_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.catalogs
    ADD CONSTRAINT catalogs_pkey PRIMARY KEY (id);


--
-- Name: controls controls_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.controls
    ADD CONSTRAINT controls_pkey PRIMARY KEY (id);


--
-- Name: metadata metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.metadata
    ADD CONSTRAINT metadata_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (id);


--
-- Name: results results_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT results_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: assessments fk_assessments_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.assessments
    ADD CONSTRAINT fk_assessments_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: catalogs fk_catalogs_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.catalogs
    ADD CONSTRAINT fk_catalogs_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: controls fk_controls_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.controls
    ADD CONSTRAINT fk_controls_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: controls fk_controls_profile_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.controls
    ADD CONSTRAINT fk_controls_profile_id FOREIGN KEY (profile_id) REFERENCES public.profiles(id);


--
-- Name: profiles fk_profiles_catalog_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT fk_profiles_catalog_id FOREIGN KEY (catalog_id) REFERENCES public.catalogs(id);


--
-- Name: profiles fk_profiles_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT fk_profiles_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: results fk_results_assessment_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT fk_results_assessment_id FOREIGN KEY (assessment_id) REFERENCES public.assessments(id);


--
-- Name: results fk_results_control_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT fk_results_control_id FOREIGN KEY (control_id) REFERENCES public.controls(id);


--
-- Name: results fk_results_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT fk_results_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: results fk_results_subject_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.results
    ADD CONSTRAINT fk_results_subject_id FOREIGN KEY (subject_id) REFERENCES public.subjects(id);


--
-- Name: subjects fk_subjects_metadata_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT fk_subjects_metadata_id FOREIGN KEY (metadata_id) REFERENCES public.metadata(id);


--
-- Name: subjects fk_subjects_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: dbadmin
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT fk_subjects_parent_id FOREIGN KEY (parent_id) REFERENCES public.subjects(id);


--
-- PostgreSQL database dump complete
--

