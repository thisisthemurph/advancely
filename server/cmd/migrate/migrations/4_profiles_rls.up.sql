alter table profiles enable row level security;

create policy "Allow select for authenticated users of the same company"
on public.profiles
for select
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.company_id = profiles.company_id
    )
);

create policy "Allow update for own profile"
on public.profiles
for update
using (auth.uid() = profiles.id);

create policy "Allow update for admin users"
on public.profiles
for update
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.is_admin = true
    )
);

create policy "Allow delete for admin users"
on public.profiles
for delete
using (
    exists (
        select 1 from public.profiles as p
        where p.id = auth.uid()
            and p.is_admin = true
    )
)
