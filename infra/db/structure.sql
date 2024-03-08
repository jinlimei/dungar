
create sequence markov_id_seq;

create table sentence
(
  id serial not null
    constraint sentence_pkey
      primary key,
  timestamp integer not null
);

create table word
(
  id serial not null
    constraint word_pkey
      primary key,
  word varchar(255) not null
);

create table fragment
(
  id serial not null
    constraint fragment_pkey
      primary key,
  r_word_id integer not null
    constraint fragment_r_word_id_fkey
      references word,
  l_word_id integer not null
    constraint fragment_l_word_id_fkey
      references word,
  sentence_id integer not null
    constraint fragment_sentence_id_fkey
      references sentence,
  word_id integer not null
    constraint fragment_word_id_fkey
      references word
);

create index idx_l_word
  on fragment (word_id, r_word_id);

create index idx_r_word
  on fragment (l_word_id, word_id);

create table idx_word_fragment
(
  id serial not null
    constraint idx_word_fragment_pkey
      primary key,
  word_id integer not null
    constraint idx_word_fragment_word_id_fkey
      references word,
  fragment_id integer not null
    constraint idx_word_fragment_fragment_id_fkey
      references fragment,
  constraint idx_word_fragment_word_id_fragment_id_key
    unique (word_id, fragment_id)
);

create index ix_idx_word_fragment_fragment_id
  on idx_word_fragment (fragment_id);

create index ix_idx_word_fragment_word_id
  on idx_word_fragment (word_id);

create unique index ix_word_word
  on word (word);

create table xmpp_jid
(
  id serial not null
    constraint xmpp_jid_pkey
      primary key,
  value varchar(255) not null
);

create unique index ix_xmpp_jid_value
  on xmpp_jid (value);

create table xmpp_muc
(
  id serial not null
    constraint xmpp_muc_pkey
      primary key,
  value varchar(255) not null
);

create unique index ix_xmpp_muc_value
  on xmpp_muc (value);

create table xmpp_nickname
(
  id serial not null
    constraint xmpp_nickname_pkey
      primary key,
  value varchar(255) not null
);

create table xmpp_log
(
  id serial not null
    constraint xmpp_log_pkey
      primary key,
  muc_id integer not null
    constraint xmpp_log_muc_id_fkey
      references xmpp_muc,
  sentence_id integer not null
    constraint xmpp_log_sentence_id_fkey
      references sentence,
  nickname_id integer not null
    constraint xmpp_log_nickname_id_fkey
      references xmpp_nickname,
  jid_id integer not null
    constraint xmpp_log_jid_id_fkey
      references xmpp_jid
);

create unique index ix_xmpp_nickname_value
  on xmpp_nickname (value);

create table gdpr_users
(
  username varchar(255) not null
    constraint gdpr_users_username_pk
      primary key,
  state integer not null,
  time integer not null,
  expires integer not null
);

create table fortunes
(
  id serial not null
    constraint fortunes_pkey
      primary key,
  fortune text,
  active smallint default 1 not null,
  last_used timestamp,
  added timestamp not null
);

create index active_idx
  on fortunes (active);

create table raw_messages
(
  id serial not null
    constraint raw_messages_pk
      primary key,
  message text not null,
  source varchar(50) not null,
  created_at timestamp with time zone not null
);

create table if not exists log_issues
(
    id serial not null
        constraint log_issues_pk
            primary key,
    issue_type varchar(50) not null,
    title varchar(150) not null,
    contents text not null,
    created timestamp not null
);

create table if not exists user_tracking
(
    id serial not null
        constraint table_name_pk
            primary key,
    unique_id varchar(20) not null,
    nick varchar(30) not null,
    line_count integer default 0 not null,
    first_seen timestamp not null,
    last_seen timestamp not null
);

create unique index if not exists table_name_unique_id_uindex
    on user_tracking (unique_id);

create table bad_words
(
    id serial not null
        constraint bad_words_pk
            primary key,
    word varchar(50) not null,
    active int default 1 not null,
    added timestamp not null
);

-- we need to prime our word table for empty stuff
-- empty stuff is represented by a id=1 being '', apparently
insert into word (word)
values ('');

CREATE SEQUENCE public.dad_jokes_id_seq;

ALTER SEQUENCE public.dad_jokes_id_seq
    OWNER TO dungarmatic;

CREATE TABLE public.dad_jokes
(
    id integer NOT NULL DEFAULT nextval('dad_jokes_id_seq'::regclass),
    the_joke text COLLATE pg_catalog."default" NOT NULL,
    joke_checksum character varying(255) COLLATE pg_catalog."default" NOT NULL,
    joke_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp NOT NULL,
    last_used timestamp ,
    CONSTRAINT dad_jokes_pkey PRIMARY KEY (id)
)
    WITH (
        OIDS = FALSE
    )
    TABLESPACE pg_default;

ALTER TABLE public.dad_jokes
    OWNER to dungarmatic;

GRANT ALL ON TABLE public.dad_jokes TO dungarmatic;
-- Index: joke_checksum_idx

-- DROP INDEX public.joke_checksum_idx;

CREATE UNIQUE INDEX joke_checksum_idx
    ON public.dad_jokes USING btree
        (joke_checksum COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: joke_id_idx

-- DROP INDEX public.joke_id_idx;

CREATE UNIQUE INDEX joke_id_idx
    ON public.dad_jokes USING btree
        (joke_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE TABLE raw_messages_m1
(
    sentence_id integer NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    PRIMARY KEY (sentence_id)
)
    WITH (
        OIDS = FALSE
    )
    TABLESPACE pg_default;

ALTER TABLE raw_messages_m1
    OWNER to dungarmatic;

create table tomes
(
    id serial not null
        constraint tomes_pk
            primary key,
    title varchar(30) not null,
    contents text not null,
    variant varchar(30) not null,
    active integer default 0 not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);

alter table tomes owner to postgres;

create unique index tomes_title_uindex
    on tomes (title);

