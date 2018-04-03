DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'source_type') THEN
    create type source_type as enum ('xml', 'csv', 'json');
  END IF;
END$$;

BEGIN;

create table if not exists hotels (
  hotel_id     varchar(32),
  source       source_type,
  source_id    varchar(100),

  name         varchar(100),
  description  text,
  stars        smallint,

  country_code varchar(2),
  country      varchar(20),
  city         varchar(100),
  address      text,
  lat          float,
  lon          float,

  primary key (hotel_id)
);

create table if not exists photos (
  url      text,
  hotel_id varchar(32) not null references hotels (hotel_id) on delete cascade
);


COMMIT;