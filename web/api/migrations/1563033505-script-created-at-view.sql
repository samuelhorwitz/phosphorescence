create or replace view scripts_view as
select
	scripts.id,
	scripts.author_id,
	scripts.forked_from_script_id,
	scripts.is_private,
	scripts.name,
	scripts.forked_from_script_version_created_at,
	script_versions.created_at
from scripts
join (
	select script_id, min(created_at) as created_at
	from script_versions
	group by script_id
) script_versions on scripts.id = script_versions.script_id
where scripts.deleted_at is null;
