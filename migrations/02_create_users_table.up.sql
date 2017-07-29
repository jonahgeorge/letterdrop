create table users (
  id serial,
  name text not null,
  email text not null,
  password_digest text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key(id),
  unique(email)
);


