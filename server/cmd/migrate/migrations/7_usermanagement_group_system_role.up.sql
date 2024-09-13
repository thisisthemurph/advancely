insert into security.permission_groups (name, description)
values
    ('User management', 'Permissions for creating, updating, and deleting users.');

insert into security.permissions (group_id, name, description)
select id, v.name, v.description
from security.permission_groups g
join (
    values
        ('create-user', 'The ability to create a new user in your organization.'),
        ('edit-user', 'The ability to edit users in your organization.'),
        ('delete-user', 'The ability to delete users in your organization.')
) as v(name, description) on true
where g.name = 'User management';

insert into security.role_permissions (role_id, permission_id)
select r.id, p.id
from security.roles r
cross join security.permissions p
where r.name = 'Admin'
on conflict do nothing;
