alter table companies enable row level security;

create policy "Allow select for authenticated users from same company"
on companies
for select
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.company_id = companies.id
    )
);

create policy "Allow update for admin users"
on companies
for update
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.company_id = companies.id
            and p.is_admin = true
    )
);

create policy "Allow delete for admin users"
on companies
for delete
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.company_id = companies.id
            and p.is_admin = true
    )
)
