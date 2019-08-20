-- IMPORTANT! These must be run as super user!
-- create extension pg_trgm;
-- create extension unaccent;
-- create extension plperl;
-- `plperl.on_init = 'use utf8; use re; package utf8; require "utf8_heavy.pl";'` needed in the config (this is already on DO managed instance)

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

-- A hashtag begins the second a hashmark is seen regardless of what preceeds it
-- and must be followed by a letter or number from across the entire Unicode
-- spectrum. Then, optionally, more letters and numbers, potentially with single
-- dashes or connectors interspersed throughout, but no trailing dashes or
-- connectors.

create function remove_hashtags(text) returns text as $$
	$_[0] =~ s/#((?:[\pL\pN]+(?:[\p{Pc}\p{Pd}][\pL\pN]+)*))//g;
	return $_[0];
$$ language plperl;

create function match_hashtags(text) returns text[][] as $$
	@matches = ($_[0] =~ /#((?:[\pL\pN]+(?:[\p{Pc}\p{Pd}][\pL\pN]+)*))/gi);
	return [@matches];
$$ language plperl;

create function validate_hashtag(text) returns boolean as $$
	return true if $_[0] =~ /^(?:[\pL\pN]+(?:[\p{Pc}\p{Pd}][\pL\pN]+)*)$/i;
	return false;
$$ language plperl;

create function remove_allowed_connectors_from_hashtag(text) returns text as $$
	$_[0] =~ s/[\p{Pc}\p{Pd}]//g;
	return $_[0];
$$ language plperl;

create function clean_hashtag(text) returns text as $$
	select remove_allowed_connectors_from_hashtag(lower(unaccent($1)))
$$ language sql;

create function validate_and_clean_hashtag(text) returns text as $$
	select case when validate_hashtag($1) then clean_hashtag($1) else null end
$$ language sql;

create function get_hashtags(text) returns text[] as $$
	select array(select distinct unnest(array_remove(array(select clean_hashtag(array_to_string(match_hashtags($1), ','))), null)))
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
		where similarity(word, original_words) >= 0.5
		order by word <-> original_words asc
	) lex;
$$ language sql;

create type search_result as (rank real, id uuid, type searchable_type, name text, description text, author_name text, likes bigint);

create function search(text) returns setof search_result as $$
	with
	phraseA as (select phraseto_tsquery('simple', unaccent($1)) as phraseA),
	phraseB as (select phraseto_tsquery('english', unaccent($1)) as phraseB),
	phraseC as (select phraseto_tsquery('simple', get_close_phrase($1)) as phraseC),
	phraseD as (select phraseto_tsquery('english', get_close_phrase($1)) as phraseD),
	phrase as (select array_to_string(array_remove(array[
		phraseA::text,
		phraseB::text,
		phraseC::text,
		phraseD::text
	], ''),'|')::tsquery as phrase from phraseA, phraseB, phraseC, phraseD)
	select rank, id, type, name, description, author_name, likes from (
		select rank, id, type, name, description, author_name, likes
		from (
			select distinct on (id) * from (
				select (((ts_rank(document, phraseA) + (ts_rank(document, phraseB) * 0.4) + (ts_rank(document, phraseC) * 0.2) + (ts_rank(document, phraseD) * 0.1))) / 4)::real as rank,
				searchables.id,
				searchables.type,
				ts_headline('english'::regconfig, name, phrase, 'StartSel = <mark>, StopSel = </mark>') as name,
				ts_headline('english'::regconfig, unmodified_description, phrase, 'StartSel = <mark>, StopSel = </mark>') as description,
				ts_headline('simple'::regconfig, author_name, phrase, 'StartSel = <mark>, StopSel = </mark>') as author_name,
				searchables.likes
				from searchables, phrase, phraseA, phraseB, phraseC, phraseD
				union
				select greatest(similarity(unaccent($1), unaccent(searchables.name)), similarity(unaccent($1), unaccent(searchables.description))) as rank,
				searchables.id,
				searchables.type,
				ts_headline('simple'::regconfig, name, phrase, 'StartSel = <mark>, StopSel = </mark>') as name,
				ts_headline('simple'::regconfig, unmodified_description, phrase, 'StartSel = <mark>, StopSel = </mark>') as description,
				searchables.author_name,
				searchables.likes
				from searchables, phrase
				where unaccent(searchables.name) ilike ('%' || unaccent($1) || '%')
				or unaccent(searchables.description) ilike ('%' || unaccent($1) || '%')
				order by rank desc, likes desc
			) results
		) searchables
		where rank > 5.96e-08 -- epsilon
	) searchables_uq
	order by rank desc, likes desc
$$ language sql;

create function search_hashtag(text) returns setof searchables as $$
	select * from searchables where tags @> array[validate_and_clean_hashtag($1)]
	order by likes desc
$$ language sql;