insert into security.roles (company_id, name, description, is_system_role)
values
    (null, 'Admin', 'A system admin user has all permission.', true);

insert into security.permission_groups (name, description)
values
    ('Permissions', 'Permissions for creating, updating, and removing permissions from other users.');

insert into security.permissions (group_id, name, description)
select id, v.name, v.description
from security.permission_groups g
join (
    values
        ('create-role', 'A user that can create new permission roles'),
        ('edit-role', 'A user that can edit existing permission roles'),
        ('delete-role', 'A user that can delete permission roles'),

        ('assign-user-role', 'A user that can assign a permissions role to a user'),
        ('remove-user-role', 'A user that can remove a permissions role from a user')
) as v(name, description) on true
where g.name = 'Permissions';

insert into security.role_permissions (role_id, permission_id)
select r.id, p.id
from security.roles r
cross join security.permissions p;

