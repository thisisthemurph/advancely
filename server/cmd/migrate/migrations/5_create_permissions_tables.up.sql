create schema if not exists security;

-- Contains roles such as 'Admin' etc
-- Users should be able to create their own roles in addition to predefined system roles
create table if not exists security.roles (
    id serial primary key,
    company_id uuid references public.companies on delete cascade,
    name text unique not null,
    description text not null,
    is_system_role boolean not null default false, -- if false, role is user created
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

create or replace trigger trg_set_updated_at_roles
    before update on security.roles
    for each row
execute function update_updated_at_timestamp();

-- A join table between auth.users and security.roles
-- This join should not be necessary for system roles where is_system_role = true
create table if not exists security.user_roles (
    user_id uuid not null references auth.users (id) on delete cascade,
    role_id integer not null references security.roles (id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp default null,

    primary key (user_id, role_id)
);

create or replace trigger trg_set_updated_at_user_roles
    before update on security.user_roles
    for each row
execute function update_updated_at_timestamp();

-- Logical groupings for permissions, such as 'User management'.
-- This table will not be updated by the users directly.
create table if not exists security.permission_groups (
    id serial primary key,
    name text unique not null,
    description text not null
);

-- Permissions such as the ability to create or delete users.
-- This table will not be updated by the users directly.
create table if not exists security.permissions (
    id serial primary key,
    group_id integer not null references security.permission_groups on delete cascade,
    name text unique not null,
    description text not null
);

-- Join table for roles and permissions.
create table if not exists security.role_permissions (
    role_id integer not null references security.roles (id) on delete cascade,
    permission_id integer not null references security.permissions (id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp default null,

    primary key (role_id, permission_id)
);

create or replace trigger trg_set_updated_at_role_permissions
    before update on security.role_permissions
    for each row
execute function update_updated_at_timestamp();
