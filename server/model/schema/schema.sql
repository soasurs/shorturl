CREATE TABLE IF NOT EXISTS shorturl (
  id serial primary key,
  created_at timestamp with time zone not null,
  deleted_at timestamp with time zone,
  url text not null unique,
  key text not null unique
);