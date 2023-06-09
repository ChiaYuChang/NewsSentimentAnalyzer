--
-- PostgreSQL database dump
--

-- Dumped from database version 14.6 (Debian 14.6-1.pgdg110+1)
-- Dumped by pg_dump version 15.3

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: api_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.api_type AS ENUM (
    'language_model',
    'source'
);


ALTER TYPE public.api_type OWNER TO postgres;

--
-- Name: event_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.event_type AS ENUM (
    'sign-in',
    'sign-out',
    'authorization',
    'api-key',
    'query'
);


ALTER TYPE public.event_type OWNER TO postgres;

--
-- Name: job_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.job_status AS ENUM (
    'created',
    'running',
    'done',
    'failure',
    'canceled'
);


ALTER TYPE public.job_status OWNER TO postgres;

--
-- Name: role; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.role AS ENUM (
    'user',
    'admin'
);


ALTER TYPE public.role OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apikeys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.apikeys (
    id integer NOT NULL,
    owner integer NOT NULL,
    api_id smallint NOT NULL,
    key text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.apikeys OWNER TO postgres;

--
-- Name: apikeys_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.apikeys_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apikeys_id_seq OWNER TO postgres;

--
-- Name: apikeys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.apikeys_id_seq OWNED BY public.apikeys.id;


--
-- Name: apis; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.apis (
    id smallint NOT NULL,
    name character varying(20) NOT NULL,
    type public.api_type NOT NULL,
    image character varying(128) DEFAULT 'logo_Default.svg'::character varying NOT NULL,
    icon character varying(128) DEFAULT 'static/image/icon_Default.svg'::character varying NOT NULL,
    document_url character varying(128) DEFAULT '#'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.apis OWNER TO postgres;

--
-- Name: apis_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.apis_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apis_id_seq OWNER TO postgres;

--
-- Name: apis_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.apis_id_seq OWNED BY public.apis.id;


--
-- Name: endpoints; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.endpoints (
    id integer NOT NULL,
    name character varying(32) NOT NULL,
    api_id smallint NOT NULL,
    template_name character varying(32) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.endpoints OWNER TO postgres;

--
-- Name: endpoints_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.endpoints_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.endpoints_id_seq OWNER TO postgres;

--
-- Name: endpoints_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.endpoints_id_seq OWNED BY public.endpoints.id;


--
-- Name: jobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.jobs (
    id integer NOT NULL,
    owner integer NOT NULL,
    status public.job_status NOT NULL,
    src_api_id smallint NOT NULL,
    src_query text NOT NULL,
    llm_api_id smallint NOT NULL,
    llm_query json NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.jobs OWNER TO postgres;

--
-- Name: keywords; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.keywords (
    id bigint NOT NULL,
    news_id bigint NOT NULL,
    keyword character varying(50) NOT NULL
);


ALTER TABLE public.keywords OWNER TO postgres;

--
-- Name: keywords_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.keywords_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.keywords_id_seq OWNER TO postgres;

--
-- Name: keywords_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.keywords_id_seq OWNED BY public.keywords.id;


--
-- Name: logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.logs (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    type public.event_type NOT NULL,
    message character varying(256) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.logs OWNER TO postgres;

--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.logs_id_seq OWNER TO postgres;

--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


--
-- Name: news; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.news (
    id bigint NOT NULL,
    md5_hash character(128) NOT NULL,
    title text NOT NULL,
    url text NOT NULL,
    description text NOT NULL,
    content text NOT NULL,
    source text,
    publish_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.news OWNER TO postgres;

--
-- Name: news_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.news_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.news_id_seq OWNER TO postgres;

--
-- Name: news_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.news_id_seq OWNED BY public.news.id;


--
-- Name: newsjobs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.newsjobs (
    id bigint NOT NULL,
    job_id bigint NOT NULL,
    news_id bigint NOT NULL
);


ALTER TABLE public.newsjobs OWNER TO postgres;

--
-- Name: newsjobs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.newsjobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.newsjobs_id_seq OWNER TO postgres;

--
-- Name: newsjobs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.newsjobs_id_seq OWNED BY public.newsjobs.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    password bytea NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    role public.role NOT NULL,
    email character varying(320) NOT NULL,
    opt character varying(128) DEFAULT NULL::character varying,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    password_updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: apikeys id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apikeys ALTER COLUMN id SET DEFAULT nextval('public.apikeys_id_seq'::regclass);


--
-- Name: apis id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apis ALTER COLUMN id SET DEFAULT nextval('public.apis_id_seq'::regclass);


--
-- Name: endpoints id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.endpoints ALTER COLUMN id SET DEFAULT nextval('public.endpoints_id_seq'::regclass);


--
-- Name: keywords id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keywords ALTER COLUMN id SET DEFAULT nextval('public.keywords_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: news id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news ALTER COLUMN id SET DEFAULT nextval('public.news_id_seq'::regclass);


--
-- Name: newsjobs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newsjobs ALTER COLUMN id SET DEFAULT nextval('public.newsjobs_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: apikeys apikeys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apikeys
    ADD CONSTRAINT apikeys_pkey PRIMARY KEY (id);


--
-- Name: apis apis_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apis
    ADD CONSTRAINT apis_pkey PRIMARY KEY (id);


--
-- Name: endpoints endpoints_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.endpoints
    ADD CONSTRAINT endpoints_pkey PRIMARY KEY (id);


--
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--
-- Name: keywords keywords_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keywords
    ADD CONSTRAINT keywords_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: news news_md5_hash_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news
    ADD CONSTRAINT news_md5_hash_key UNIQUE (md5_hash);


--
-- Name: news news_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news
    ADD CONSTRAINT news_pkey PRIMARY KEY (id);


--
-- Name: newsjobs newsjobs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newsjobs
    ADD CONSTRAINT newsjobs_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: apikeys_owner_api_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX apikeys_owner_api_id_idx ON public.apikeys USING btree (owner, api_id);


--
-- Name: jobs_owner_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX jobs_owner_status_idx ON public.jobs USING btree (owner, status);


--
-- Name: keywords_keyword_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX keywords_keyword_idx ON public.keywords USING btree (keyword);


--
-- Name: logs_user_id_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX logs_user_id_type_idx ON public.logs USING btree (user_id, type);


--
-- Name: news_md5_hash_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX news_md5_hash_idx ON public.news USING btree (md5_hash);


--
-- Name: news_publish_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX news_publish_at_idx ON public.news USING btree (publish_at);


--
-- Name: newsjobs_job_id_news_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX newsjobs_job_id_news_id_idx ON public.newsjobs USING btree (job_id, news_id);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: apikeys apikeys_api_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apikeys
    ADD CONSTRAINT apikeys_api_id_fkey FOREIGN KEY (api_id) REFERENCES public.apis(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: apikeys apikeys_owner_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apikeys
    ADD CONSTRAINT apikeys_owner_fkey FOREIGN KEY (owner) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: endpoints endpoints_api_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.endpoints
    ADD CONSTRAINT endpoints_api_id_fkey FOREIGN KEY (api_id) REFERENCES public.apis(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: jobs jobs_llm_api_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_llm_api_id_fkey FOREIGN KEY (llm_api_id) REFERENCES public.apis(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: jobs jobs_owner_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_owner_fkey FOREIGN KEY (owner) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: jobs jobs_src_api_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_src_api_id_fkey FOREIGN KEY (src_api_id) REFERENCES public.apis(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: keywords keywords_news_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.keywords
    ADD CONSTRAINT keywords_news_id_fkey FOREIGN KEY (news_id) REFERENCES public.news(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: logs logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE;


--
-- Name: newsjobs newsjobs_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newsjobs
    ADD CONSTRAINT newsjobs_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: newsjobs newsjobs_news_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newsjobs
    ADD CONSTRAINT newsjobs_news_id_fkey FOREIGN KEY (news_id) REFERENCES public.news(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

