# Advancely

## Local setup

When running locally in development mode, the React front-end and the Go backend server are ran separately.

In production, the go backend serves the React front-end.

### Emails with Resend

We use the [Resend](https://resend.com/) emailing service.

Request an API key for later configuration in Supabase.

### Configure Supabase

We use [Supabase](https://supabase.com/) for the Postgres database and authentication.

Go to [the Supabase website](https://supabase.com/), create an account and a new project named something like `Advancely-dev`.

Make the following configuration changes:

**`Authentication > URL Configuration`**

- Set the `Site URL` to `http://localhost:5173`.
- Add the following `Redirect URLs`:
  - http://localhost:5173/auth/callback

**`Settings > Authentication :: SMTP Settings`**

- Ensure `Enable Custom SMTP` is enabled
- Set the following values:
  - Sender email: `you@toptabletrivia.com`
  - Sender name: `a generic name`
  - Host: `smtp.resend.com`
  - Port number: `465`
  - Username: `resend`
  - Password: `your resend API key`

### Environment variables

Ensure the following environment variables are set:

**client/.env**

```properties
VITE_API_BASE_URL=http://localhost:42069/api/v1
```

**server/.env**

```properties
ENVIRONMENT=development
LISTEN_ADDRESS=:42069
CLIENT_BASE_URL=http://localhost:5173
SESSION_SECRET=session.secret

# This information can be obtained from your Supabase settings
# Navigate to `Settings > Database`
DATABASE_NAME=
DATABASE_PASSWORD=
DATABASE_URI=
AUTO_MIGRATE_ON=false

# This information can be obtained from your Supabase settings
# Navigate to `Settings > API`
SUPABASE_URL=
SUPABASE_PUBLIC_KEY=
SUPABASE_SERVICE_ROLE_SECRET=

RESEND_KEY=<create a resend key at resend.com>
```

## Running locally

Run the server; this will run commands present in the `server/Makefile`:

```console
make api
```

Install node dependencies:

```console
npm install
```

If you want to watch for tailwind changes, you should run:

```console
npm run tailwind
```

Run the front-end react app:

```console
npm run dev
```

## Database migrations

When deploying in production, the database will automatically be migrated up.

In development, this will only happen if the `AUTO_MIGRATE_ON` environment variable is set to `true`.

To run a migration manually you can run the following from the `server` directory:

```sh
make migrate
```

To run a down migration, you can run the following:

```sh
make migrate-down
```

Note that running the down migration will clear all data from the database and drop all tables.
