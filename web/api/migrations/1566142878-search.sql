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
	description,
	tags,
	author_name,
	build_document(name, description, author_name, tags, true) as document
from (
	select
		scripts.id as id,
		'script'::searchable_type as type,
		scripts.name as name,
		remove_hashtags(scripts.description) as description,
		get_hashtags(scripts.description) as tags,
		users.name as author_name
	from scripts_public_view scripts
	left join users_view users on users.id = scripts.author_id
) scripts
union
select
	id,
	type,
	name,
	description,
	tags,
	author_name,
	build_document(name, description, author_name, tags, true) as document
from (
	select
		script_chains.id as id,
		'script_chain'::searchable_type as type,
		script_chains.name as name,
		remove_hashtags(script_chains.description) as description,
		get_hashtags(script_chains.description) as tags,
		users.name as author_name
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
