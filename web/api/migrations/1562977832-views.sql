-- These are for convenience not security, so no `security_barrier` option is being used
create view users_view as select id, spotify_id from users where deleted_at is null;
create view scripts_view as select id, author_id, forked_from_script_id, is_private from scripts where deleted_at is null;
create view script_versions_view as select script_id, created_at, type, file_id from script_versions where deleted_at is null;

grant select on users_view to phosphor_api;
grant select on scripts_view to phosphor_api;
grant select on script_versions_view to phosphor_api;
