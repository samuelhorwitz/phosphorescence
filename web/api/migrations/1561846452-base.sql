do $$ begin
    if not exists (select from pg_catalog.pg_roles where rolname = 'phosphor_api') then
        create role phosphor_api login password 'phosphor_api'; -- this is for local use only, role exists in production with better password
    end if;
end $$;

create table _migrations (
    id uuid primary key,
    name text not null,
    ran_at timestamp with time zone not null default now()
);

create table users (
    id uuid primary key,
    spotify_id text not null unique, -- spotify IDs are integers but best practices for IDs we don't own are to treat them as opaque strings
    deleted_at timestamp with time zone
);

create table scripts (
    id uuid primary key,
    author_id uuid not null references users(id) on update restrict on delete restrict,
    forked_from_script_id uuid references scripts(id) on update restrict on delete restrict,
    is_private boolean not null default true,
    deleted_at timestamp with time zone
);

create index on scripts (author_id);
create index on scripts (forked_from_script_id);

create type version_type as enum (
    'draft',
    'autosave',
    'publish'
);

create table script_versions (
    script_id uuid not null references scripts(id) on update restrict on delete restrict,
    created_at timestamp with time zone not null default now(),
    deleted_at timestamp with time zone,
    type version_type not null,
    file_id uuid not null, -- actual content is stored in object storage
    primary key (script_id, created_at) -- cannot have two versions from exact same time, simplifies schema and ordering and is realistically non-issue
);

create index on script_versions (file_id);

grant connect on database phosphor to phosphor_api;
grant usage on schema public to phosphor_api;
grant select on users to phosphor_api;
grant insert on users to phosphor_api;
grant update (spotify_id, deleted_at) on users to phosphor_api;
grant select on scripts to phosphor_api;
grant insert on scripts to phosphor_api;
grant update (is_private, deleted_at) on scripts to phosphor_api;
grant select on script_versions to phosphor_api;
grant insert on script_versions to phosphor_api;
grant update (deleted_at) on script_versions to phosphor_api;
