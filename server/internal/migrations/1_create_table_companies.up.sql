create table if not exists companies (
    id uuid primary key default uuid_generate_v4(),
    name text not null,
    creator_id uuid references auth.users (id) on delete set null,
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

create trigger set_updated_at
    before update on companies
    for each row
        execute function update_updated_at_timestamp();
