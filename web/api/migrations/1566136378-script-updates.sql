alter table scripts drop constraint scripts_forked_from_script_id_fkey;
alter table scripts add column description text;
grant update (name, description) on scripts to phosphor_api;

create table script_likes (
    script_id uuid not null references scripts(id) on update restrict on delete restrict,
    liker_id uuid not null references users(id) on update restrict on delete restrict,
    primary key (script_id, liker_id)
);

create index on script_likes (liker_id);

grant insert on script_likes to phosphor_api;
grant delete on script_likes to phosphor_api;

create table script_chains (
	id uuid primary key,
	author_id uuid not null references users(id) on update restrict on delete restrict,
	forked_from_script_chain_id uuid,
	is_private boolean not null default true,
	deleted_at timestamp with time zone,
	name text,
	forked_from_script_chain_version_created_at timestamp with time zone,
	description text
);

alter table script_chains add check
	((forked_from_script_chain_id is null and forked_from_script_chain_version_created_at is null) or
	(forked_from_script_chain_id is not null and forked_from_script_chain_version_created_at is not null));

create index on script_chains (author_id);
create index on script_chains (forked_from_script_chain_id);

create table script_chain_versions (
    script_chain_id uuid not null references script_chains(id) on update restrict on delete restrict,
    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone,
    seeder_id uuid references scripts(id) on update restrict on delete restrict,
    builder_id uuid not null references scripts(id) on update restrict on delete restrict,
    primary key (script_chain_id, created_at) -- cannot have two versions from exact same time, simplifies schema and ordering and is realistically non-issue
);

create index on script_chain_versions (seeder_id);
create index on script_chain_versions (builder_id);

alter table script_chains add foreign key (forked_from_script_chain_id, forked_from_script_chain_version_created_at)
	references script_chain_versions(script_chain_id, created_at) on update restrict on delete restrict;

create table script_chain_pruners (
	script_chain_id uuid not null,
	created_at timestamp with time zone not null,
	pruner_id uuid not null references scripts(id) on update restrict on delete restrict,
	execution_order smallint not null,
	primary key (script_chain_id, created_at, pruner_id)
);

create index on script_chain_pruners (pruner_id);

alter table script_chain_pruners add foreign key (script_chain_id, created_at)
	references script_chain_versions(script_chain_id, created_at) on update restrict on delete restrict;

grant select on script_chains to phosphor_api;
grant insert on script_chains to phosphor_api;
grant update (is_private, deleted_at, name, description) on script_chains to phosphor_api;
grant select on script_chain_versions to phosphor_api;
grant insert on script_chain_versions to phosphor_api;
grant update (deleted_at) on script_chain_versions to phosphor_api;
grant select on script_chain_pruners to phosphor_api;
grant insert on script_chain_pruners to phosphor_api;

create or replace view scripts_view as
select
	scripts.id,
	scripts.author_id,
	scripts.forked_from_script_id,
	scripts.is_private,
	scripts.name,
	scripts.forked_from_script_version_created_at,
	script_versions.created_at,
	scripts.description
from scripts
join (
	select script_id, min(created_at) as created_at
	from script_versions_view
	group by script_id
) script_versions on scripts.id = script_versions.script_id
where scripts.deleted_at is null;

create view script_chain_versions_view as
select
	script_chain_versions.script_chain_id,
	script_chain_versions.created_at,
	script_chain_versions.seeder_id,
	script_chain_versions.builder_id,
	script_chain_pruners.pruner_ids
from script_chain_versions
left join (
	select script_chain_id, created_at, array_agg(pruner_id order by execution_order asc) as pruner_ids
	from script_chain_pruners
	group by script_chain_id, created_at
) script_chain_pruners on script_chain_versions.script_chain_id = script_chain_pruners.script_chain_id
	and script_chain_versions.created_at = script_chain_pruners.created_at
where script_chain_versions.deleted_at is null;

create view script_chains_view as
select
	script_chains.id,
	script_chains.author_id,
	script_chains.forked_from_script_chain_id,
	script_chains.is_private,
	script_chains.name,
	script_chains.forked_from_script_chain_version_created_at,
	script_chain_versions.created_at,
	script_chains.description
from script_chains
join (
	select script_chain_id, min(created_at) as created_at
	from script_chain_versions_view
	group by script_chain_id
) script_chain_versions on script_chains.id = script_chain_versions.script_chain_id
where script_chains.deleted_at is null;
