letterbox
===

```sql
create extension pgcrypto;

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
  name text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key (id)
);
```
