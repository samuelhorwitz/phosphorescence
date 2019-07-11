alter table scripts add column name text; -- null is "untitled"
-- created_at is actually the version id, incremental but not abstract
alter table scripts add column forked_from_script_version_created_at timestamp with time zone;
alter table scripts add foreign key (forked_from_script_id, forked_from_script_version_created_at)
	references script_versions(script_id, created_at) on update restrict on delete restrict;
alter table scripts add check
	((forked_from_script_id is null and forked_from_script_version_created_at is null) or
	(forked_from_script_id is not null and forked_from_script_version_created_at is not null));

create or replace view scripts_view as select
	id, author_id, forked_from_script_id, is_private, name, forked_from_script_version_created_at
from scripts where deleted_at is null;
