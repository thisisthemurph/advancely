create table if not exists profiles (
    id uuid references auth.users (id) on delete cascade,
    company_id uuid references companies (id) on delete cascade,
    first_name text not null,
    last_name text not null,
    is_admin boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp default null,

    primary key (id)
);

create trigger trg_set_updated_at_profiles
    before update on profiles
    for each row
        execute function update_updated_at_timestamp();
