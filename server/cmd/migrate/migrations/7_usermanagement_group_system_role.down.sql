with rp_to_delete as (
    select rp.role_id, rp.permission_id
    from security.role_permissions rp
    join security.permissions p
      on rp.permission_id = p.id
    join security.permission_groups g
      on p.group_id = g.id
    where g.name = 'User management'
        and role_id = (select id from security.roles where name = 'Admin')
)
delete from security.role_permissions
using rp_to_delete
where security.role_permissions.role_id = rp_to_delete.role_id
and security.role_permissions.permission_id = rp_to_delete.permission_id;

with permissions_to_delete as (
    select p.id
    from security.role_permissions rp
    join security.permissions p
        on rp.permission_id = p.id
    join security.permission_groups g
        on p.group_id = g.id
    where g.name = 'User management'
        and role_id = (select id from security.roles where name = 'Admin')
)
delete from security.permissions
using permissions_to_delete
where security.permissions.id = permissions_to_delete.id;

delete from security.permission_groups
where name = 'User management';
