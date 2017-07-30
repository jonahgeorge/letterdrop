alter table users add column email_confirmation_token text;
alter table users add column is_email_confirmed bool not null default false;
