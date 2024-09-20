create table if not exists allowed_email_domains (
    id serial primary key,
    company_id uuid references public.companies on delete cascade,
    domain text not null,
    created_at timestamp not null default now(),
    updated_at timestamp default null
);

create or replace trigger trg_set_updated_at_allowed_email_domains
    before update on public.allowed_email_domains
    for each row
execute function update_updated_at_timestamp();

alter table allowed_email_domains
drop constraint if exists unique_company_id_domain;

alter table allowed_email_domains
add constraint unique_company_id_domain
unique (company_id, domain);