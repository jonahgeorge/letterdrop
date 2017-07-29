create table submissions (
  id serial,
  form_id serial not null references forms on delete cascade,
  body text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  primary key(id)
);

