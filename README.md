letterdrop
===

```sql
create extension "pgcrypto";
create extension "uuid-ossp";

create table users (
  id serial,
  email text not null,
  password_digest text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key (id)
);

create table forms (
  id serial,
  user_id int not null,
  uuid uuid not null default uuid_generate_v4(),
  name text not null,
  description text,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key (id)
);

create table submissions (
  id serial,
  form_id int not null,
  body text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key (id)
);
```
