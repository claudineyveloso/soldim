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
  deposit_id     BIGINT NOT NULL,
  product_id     BIGINT NOT NULL,
  saldo_fisico   INTEGER NOT NULL DEFAULT 0,
  saldo_virtual  INTEGER NOT NULL DEFAULT 0,
  created_at     TIMESTAMP NOT NULL,
  updated_at     TIMESTAMP NOT NULL
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
  product_id          BIGINT NOT NULL,
  saldo_fisico_total  INT NOT NULL,
  saldo_virtual_total INT NOT NULL,
  created_at          TIMESTAMP NOT NULL,
  updated_at          TIMESTAMP NOT NULL
);
ALTER TABLE
  "stocks"
ADD
  FOREIGN KEY ("product_id") REFERENCES "products" ("id");

DROP TABLE IF EXISTS "supplier_products";
CREATE TABLE IF NOT EXISTS supplier_products (
  id            BIGINT PRIMARY KEY,
  descricao     VARCHAR(255) NOT NULL DEFAULT '',
  codigo        BIGINT NOT NULL DEFAULT 0,
  preco_custo    FLOAT NOT NULL DEFAULT 0,
  preco_compra   FLOAT NOT NULL DEFAULT 0,
  padrao        BOOLEAN NOT NULL DEFAULT TRUE,
  supplier_id   BIGINT NOT NULL DEFAULT 0,
  product_id    BIGINT NOT NULL DEFAULT 0,
  created_at    TIMESTAMP NOT NULL,
  updated_at    TIMESTAMP NOT NULL
);
ALTER TABLE
  "supplier_products"
ADD
   FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE
  "supplier_products"
ADD
    FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id");

ALTER TABLE
  "supplier_products"
ADD
   FOREIGN KEY ("product_id") REFERENCES "products" ("id");

DROP TABLE IF EXISTS "situations";
CREATE TABLE IF NOT EXISTS situations (
  id            BIGINT PRIMARY KEY,
  descricao     VARCHAR(255) NOT NULL DEFAULT '',
  created_at    TIMESTAMP NOT NULL,
  updated_at    TIMESTAMP NOT NULL
);

DROP TABLE IF EXISTS "stores";
CREATE TABLE IF NOT EXISTS stores (
  id            BIGINT PRIMARY KEY,
  descricao     VARCHAR(255) NOT NULL DEFAULT '',
  created_at    TIMESTAMP NOT NULL,
  updated_at    TIMESTAMP NOT NULL
);


DROP TABLE IF EXISTS "sales_orders";
CREATE TABLE IF NOT EXISTS sales_orders (
  id                  BIGINT PRIMARY KEY,
  numero              INT NOT NULL,
  numeroLoja          VARCHAR(100) NOT NULL,
  data                DATE NOT NULL,
  dataSaida           DATE NOT NULL,
  dataPrevista        DATE NOT NULL,
  totalProdutos       FLOAT NOT NULL DEFAULT 0,
  totalDescontos      FLOAT NOT NULL DEFAULT 0,
  situation_id        BIGINT NOT NULL DEFAULT 0,
  store_id            BIGINT NOT NULL,
  created_at          TIMESTAMP NOT NULL,
  updated_at          TIMESTAMP NOT NULL
);

ALTER TABLE
  "sales_orders"
ADD
   FOREIGN KEY ("situation_id") REFERENCES "situations" ("id");

ALTER TABLE
  "sales_orders"
ADD
   FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

DROP TABLE IF EXISTS "products_sales_orders";
CREATE TABLE IF NOT EXISTS products_sales_orders (
  sales_order_id  BIGINT NOT NULL,
  product_id      BIGINT NOT NULL,
  quantidade      INT NOT NULL DEFAULT 0,
  created_at      TIMESTAMP NOT NULL,
  updated_at      TIMESTAMP NOT NULL
);

ALTER TABLE
  "products_sales_orders"
ADD
   FOREIGN KEY ("sales_order_id") REFERENCES "sales_orders" ("id");

ALTER TABLE
  "products_sales_orders"
ADD
   FOREIGN KEY ("product_id") REFERENCES "products" ("id");


DROP TABLE IF EXISTS "triages";
CREATE TABLE IF NOT EXISTS triages (
  id                  UUID PRIMARY KEY,
  type                VARCHAR(100) NOT NULL,
  grid                VARCHAR(100) NOT NULL,
  sku_sap             INTEGER NOT NULL,
  sku_wms             VARCHAR(100) NOT NULL,
  description         VARCHAR NOT NULL,
  cust_id             BIGINT NOT NULL,
  seller              VARCHAR(255) NOT NULL,
  quantity_supplied   INTEGER NOT NULL,
  final_quantity      INTEGER NOT NULL,
  unitary_value       FLOAT NOT NULL,
  total_value_offered FLOAT NOT NULL,
  final_total_value   FLOAT NOT NULL,
  category            VARCHAR(255) NOT NULL,
  sub_category        VARCHAR(255) NOT NULL,
  sent_to_batch       BOOLEAN NOT NULL DEFAULT FALSE,
  sent_to_bling       BOOLEAN NOT NULL DEFAULT FALSE,
  defect              BOOLEAN NOT NULL DEFAULT FALSE,
  created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create unique index sku_sap_idx on triages (sku_sap);
create unique index sku_wms_idx on triages (sku_wms);


DROP TABLE IF EXISTS "contacts";
CREATE TABLE contacts (
  id              BIGINT PRIMARY KEY,
  nome            VARCHAR(100) NOT NULL DEFAULT '',
  codigo          VARCHAR(100) NOT NULL DEFAULT '',
  situacao        VARCHAR(100) NOT NULL DEFAULT '',
  numeroDocumento VARCHAR(100) NOT NULL DEFAULT '',
  telefone        VARCHAR(100) NOT NULL DEFAULT '',
  celular         VARCHAR(100) NOT NULL DEFAULT '',
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create unique index name_contact_idx on triages (name);

DROP TABLE IF EXISTS "suppliers_users";
CREATE TABLE suppliers_users (
  supplier_id BIGINT NOT NULL DEFAULT 0,
  user_id     UUID NOT NULL DEFAULT 0,
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS "tokens";
CREATE TABLE IF NOT EXISTS tokens (
  id            UUID PRIMARY KEY,
  access_token  VARCHAR(255) NOT NULL,
  expires_in    INTEGER NOT NULL,
  token_type    VARCHAR(100) NOT NULL,
  scope         VARCHAR NOT NULL,
  refresh_token VARCHAR(255) NOT NULL
);
