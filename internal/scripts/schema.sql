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
  image_url   varchar not null,
  description varchar(255) not null,
  source      varchar(100) not null,
  price       float not null,
  promotion   boolean not null default false,
  link        varchar not null,
  search_id   UUID not null,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  CONSTRAINT fk_search
    FOREIGN KEY (search_id) 
    REFERENCES searches (id)
    ON DELETE CASCADE
);

 -- Drafts
DROP TABLE IF EXISTS "drafts";
CREATE TABLE IF NOT EXISTS drafts (
  id          UUID PRIMARY KEY,
  image_url   varchar not null,
  description varchar(255) not null,
  source      varchar(100) not null,
  price       float not null,
  promotion   boolean not null default false,
  link        varchar not null,
  search_id   UUID not null,
  created_at  timestamp not null,
  updated_at  timestamp not null
);

ALTER TABLE
   "drafts"
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
  id                          bigint PRIMARY KEY,
  idProdutoPai                bigint not null default 0,
  nome                        varchar(255) not null,
  codigo                      varchar(100) not null default '',
  preco                       float not null default 0.00,
  tipo                        varchar(10) not null default '',
  situacao                    varchar(10) not null default '',
  formato                     varchar(10) not null default '',
  descricao_curta             varchar not null default '',
  imagem_url                  varchar not null default '',
  dataValidade                date not null default '2000-01-01',
  unidade                     varchar(10) not null default '',
  pesoLiquido                 float not null default 0,
  pesoBruto                   float not null default 0,
  volumes                     int not null default 0,
  itensPorCaixa               int not null default 0,
  gtin                        varchar(100) not null default '',
  gtinEmbalagem               varchar(100) not null default '',
  tipoProducao                varchar(10) not null default '',
  condicao                    int not null default 0,
  freteGratis                 boolean not null default false,
  marca                       varchar(10) not null default '',
  descricaoComplementar       varchar not null default '',
  linkExterno                 varchar not null default '',
  observacoes                 varchar not null default '',
  descricaoEmbalagemDiscreta  varchar not null default '',
  created_at                  timestamp not null,
  updated_at                  timestamp not null
);


DROP TABLE IF EXISTS "sales_channel";
CREATE TABLE sales_channel (
  id            int PRIMARY KEY,
  descricao     varchar not null,
  tipo          varchar(100) not null,
  situacao      int not null default 0,
  created_at    timestamp not null,
  updated_at    timestamp not null
);


DROP TABLE IF EXISTS "suppliers";
CREATE TABLE IF NOT EXISTS suppliers (
  id            bigint PRIMARY KEY,
  descricao     varchar(255) not null,
  codigo        varchar(255) not null,
  precoCusto    varchar(20) not null,
  precoCompra   varchar(20) not null,
  padrão        boolean not null default false,
  product_id    bigint not null,
  created_at    timestamp not null,
  updated_at    timestamp not null
);
ALTER TABLE
  "suppliers"
ADD
   FOREIGN KEY ("product_id") REFERENCES "products" ("id"); 

DROP TABLE IF EXISTS "deposits";
CREATE TABLE IF NOT EXISTS deposits (
  id                  bigint PRIMARY KEY,
  descricao           varchar(255) not null,
  situacao            int not null default 0,
  padrao              boolean not null default false,
  desconsiderarSaldo  boolean not null default false,
  created_at          timestamp not null,
  updated_at          timestamp not null
);

DROP TABLE IF EXISTS "deposit_products";
CREATE TABLE IF NOT EXISTS deposit_products (
  id            bigint PRIMARY KEY,
  deposit_id    bigint not null,
  product_id    bigint not null,
  saldoFisico   int not null default 0,
  saldoVirtual  int not null default 0,
  created_at    timestamp not null,
  updated_at    timestamp not null
);

ALTER TABLE
  "deposit_products"
ADD
    FOREIGN KEY ("deposit_id") REFERENCES "deposits" ("id"); 

ALTER TABLE
  "deposit_products"
ADD
   FOREIGN KEY ("product_id") REFERENCES "products" ("id"); 

DROP TABLE IF EXISTS "stocks";
CREATE TABLE IF NOT EXISTS stocks (
  id                bigint PRIMARY KEY,
  product_id        bigint not null,
  saldoFisicoTotal  int not null default 0,
  saldoVirtualTotal int not null default 0,
  created_at        timestamp not null,
  updated_at        timestamp not null
);
ALTER TABLE
  "stocks"
ADD
  FOREIGN KEY ("product_id") REFERENCES "products" ("id"); 

DROP TABLE IF EXISTS "tokens";
CREATE TABLE IF NOT EXISTS tokens (
  id            UUID PRIMARY KEY,
  access_token  varchar(255) not null,
  expires_in    integer not null,
  token_type    varchar(100) not null,
  scope         varchar not null,
  refresh_token varchar(255) not null
);
