-- Users
DROP TABLE IF EXISTS "users";
CREATE TABLE IF NOT EXISTS users (
  id          UUID PRIMARY KEY,
  email       varchar(100) not null,
  password    varchar(100) not null,
  is_active   boolean not null default true,
  user_type   varchar(50) not null,
  created_at  timestamp not null,
  updated_at  timestamp not null
);
create unique index email_idx on users (email);

-- Person
CREATE TABLE IF NOT EXISTS people (
  ID              UUID PRIMARY KEY,
  first_name      varchar(100) not null,
  last_name       varchar(100) not null,
  email           varchar(255) not null,
  phone           varchar(20) DEFAULT '',
  cell_phone      varchar(20) not null,
  personable_id   UUID not null,
  personable_type varchar(255) not null,
  created_at      timestamp not null,
  updated_at      timestamp not null
);
create unique index first_name_idx on people (first_name);
create unique index last_name_idx on people (last_name);
create unique index email_people_idx on people (email);

-- Address
CREATE TABLE IF NOT EXISTS addresses (
  ID                UUID PRIMARY KEY,
  public_place      varchar(255) default '',
  complement        varchar(255) default '',
  neighborhood      varchar(255) default '',
  city              varchar(255) default '',
  state             varchar(255) default '',
  zip_code          varchar(255) default '',
  addressable_id    UUID NOT NULL,
  addressable_type  varchar(255) not null,
  created_at        timestamp not null,
  updated_at        timestamp not null
);

-- Searches
DROP TABLE IF EXISTS "searches";
CREATE TABLE IF NOT EXISTS searches (
  id          UUID PRIMARY KEY,
  description varchar(255) not null,
  created_at  timestamp not null,
  updated_at  timestamp not null
);
create unique index description_idx on searches (description);

-- Searches_result
DROP TABLE IF EXISTS "searches_result";
CREATE TABLE IF NOT EXISTS searches_result (
  id          UUID PRIMARY KEY,
  image_url   varchar(255) not null,
  description varchar(255) not null,
  source      varchar(100) not null,
  price       float not null,
  promotion   boolean not null default false,
  link        varchar(255) not null,
  search_id   UUID not null,
  created_at  timestamp not null,
  updated_at  timestamp not null
);

ALTER TABLE
   "searches_result"
ADD
   FOREIGN KEY ("search_id") REFERENCES "searches" ("id");

DROP TABLE IF EXISTS "parameters";
CREATE TABLE IF NOT EXISTS parameters (
  id                  UUID PRIMARY KEY,
  discount_percentage integer NOT null default 0,
  created_at          timestamp not null,
  updated_at          timestamp not null
);
DROP TABLE IF EXISTS "products";
CREATE TABLE IF NOT EXISTS products (
  id              integer PRIMARY KEY,
  nome            varchar(255) not null,
  codigo          varchar(100) not null default '',
  preco           float not null default 0.00,
  tipo            varchar(10) not null default '',
  situacao        varchar(10) not null default '',
  formato         varchar(10) not null default '',
  descricaoCurta  varchar(10) not null default '',
  imagemURL       varchar(10) not null default '',
  created_at      timestamp not null,
  updated_at      timestamp not null
);

DROP TABLE IF EXISTS "tokens";
CREATE TABLE IF NOT EXISTS tokens (
  id            UUID PRIMARY KEY,
  access_token  varchar(255) not null,
  expires_in    integer not null,
  token_type    varchar(100) not null,
  scope         varchar(255) not null,
  refresh_token varchar(255) not null,
  created_at    timestamp not null,
  updated_at    timestamp not null
);
