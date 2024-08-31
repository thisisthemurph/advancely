drop policy if exists "Allow delete for admin users" on profiles;
drop policy if exists "Allow update for admin users" on profiles;
drop policy if exists "Allow update for own profile" on profiles;
drop policy if exists "Allow select for authenticated users of the same company" on profiles;
alter table profiles disable row level security;
