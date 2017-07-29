create table forms (
  id serial,
  user_id serial not null references users on delete cascade,
  uuid uuid not null default uuid_generate_v4(),
  name text not null,
  description text,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key(id)
);


