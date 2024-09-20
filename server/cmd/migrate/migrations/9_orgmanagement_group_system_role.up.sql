insert into security.permission_groups (name, description)
values
    ('Organization management', 'Permissions for updating details concerning the organization as a whole.');

insert into security.permissions (group_id, name, description)
select id, v.name, v.description
from security.permission_groups g
join (
    values
        ('edit-organization-settings', 'The ability to edit settings relating to the organization as a whole.')
) as v(name, description) on true
where g.name = 'Organization management';

insert into security.role_permissions (role_id, permission_id)
select r.id, p.id
from security.roles r
cross join security.permissions p
where r.name = 'Admin'
on conflict do nothing;
