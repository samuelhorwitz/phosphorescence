-- IMPORTANT! These must be run as super user!
-- create extension pg_trgm;
-- create extension unaccent;

alter table users add column name text not null;

create or replace view users_view as select id, spotify_id, name from users where deleted_at is null;

create view scripts_public_view as select id, author_id, forked_from_script_id, name, forked_from_script_version_created_at, created_at, description
from scripts_view
where is_private = false;

create view script_chains_public_view as select id, author_id, forked_from_script_chain_id, name, forked_from_script_chain_version_created_at, created_at, description
from script_chains_view
where is_private = false;

create type searchable_type as enum (
    'script',
    'script_chain'
);

create function build_document(text, text, text, text[], bool) returns tsvector as $$
	select
	setweight(to_tsvector(case when $5 then 'english'::regconfig else 'simple'::regconfig end, unaccent($1)), 'A') ||
	setweight(to_tsvector(case when $5 then 'english'::regconfig else 'simple'::regconfig end, unaccent($2)), 'B') ||
	setweight(to_tsvector('simple', unaccent($3)), 'C') ||
	setweight(to_tsvector('simple', unaccent(array_to_string($4, ' '))), 'A')
$$ language sql;

create function remove_hashtags(text) returns text as $$
	select regexp_replace($1, '#[^\s]+', '', 'g')
$$ language sql;

create function get_hashtags(text) returns text[] as $$
	select array(select unaccent(array_to_string(regexp_matches($1, '#([^#\s]+)', 'g'), ',')))
$$ language sql;

create materialized view searchables as
select
	id,
	type,
	name,
	unmodified_description,
	description,
	tags,
	author_name,
	build_document(name, description, author_name, tags, true) as document,
	likes
from (
	select
		scripts.id as id,
		'script'::searchable_type as type,
		scripts.name as name,
		scripts.description as unmodified_description,
		remove_hashtags(scripts.description) as description,
		get_hashtags(scripts.description) as tags,
		users.name as author_name,
		(select count(*) from script_likes where scripts.id = script_likes.script_id) as likes
	from scripts_public_view scripts
	left join users_view users on users.id = scripts.author_id
) scripts
union
select
	id,
	type,
	name,
	unmodified_description,
	description,
	tags,
	author_name,
	build_document(name, description, author_name, tags, true) as document,
	likes
from (
	select
		script_chains.id as id,
		'script_chain'::searchable_type as type,
		script_chains.name as name,
		script_chains.description as unmodified_description,
		remove_hashtags(script_chains.description) as description,
		get_hashtags(script_chains.description) as tags,
		users.name as author_name,
		(select count(*) from script_chain_likes where script_chains.id = script_chain_likes.script_chain_id) as likes
	from script_chains_public_view script_chains
	left join users_view users on users.id = script_chains.author_id
) script_chains;

create unique index on searchables (id);
create index on searchables using gin(tags);
create index on searchables using gin(name gin_trgm_ops);
create index on searchables using gin(author_name gin_trgm_ops);
create index on searchables using gin(document);

create materialized view searchable_lexemes as
select word from ts_stat('select build_document(name, description, author_name, tags, false) from searchables');

create index on searchable_lexemes using gin(word gin_trgm_ops);

create function get_close_phrase(text) returns text as $$
	with original_words as (select regexp_split_to_table(unaccent($1), '\s+') as original_words)
	select coalesce(string_agg(distinct lex.word, ' '), '') as corrected
	from original_words
	cross join lateral (
		select word
		from searchable_lexemes
		where similarity(word, original_words) >= 0.3
		order by word <-> original_words asc
	) lex;
$$ language sql;

create type search_result as (rank real, id uuid, type searchable_type, name text, description text, author_name text, likes bigint);

create function search(text) returns setof search_result as $$
	with phrase as (select array_to_string(array_remove(array[
		phraseto_tsquery('english', $1)::text,
		phraseto_tsquery('english', get_close_phrase($1))::text,
		phraseto_tsquery('simple', $1)::text,
		phraseto_tsquery('simple', get_close_phrase($1))::text
	], ''),'|')::tsquery as phrase)
	select rank, id, type, name, description, author_name, likes
	from (
	    select ts_rank(document, phrase) as rank,
	    searchables.id,
	    searchables.type,
	    ts_headline('english'::regconfig, name, phrase, 'StartSel = <mark>, StopSel = </mark>') as name,
	    ts_headline('english'::regconfig, unmodified_description, phrase, 'StartSel = <mark>, StopSel = </mark>') as description,
	    ts_headline('english'::regconfig, author_name, phrase, 'StartSel = <mark>, StopSel = </mark>') as author_name,
	    searchables.likes
	    from searchables, phrase
	) searchables
	where rank > 0
	order by rank desc, likes desc;
$$ language sql;