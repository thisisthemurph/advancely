drop policy if exists "Allow delete for admin users" on companies;
drop policy if exists "Allow update for admin users" on companies;
drop policy if exists "Allow select for authenticated users from same company" on companies;
alter table companies disable row level security;
