\connect stockprice

CREATE TABLE stock_data (
  "id" serial NOT NULL,
  "symbol" character varying(16) NOT NULL,
  "price" numeric(19,3) NOT NULL,
  "date" timestamp NOT NULL
);

\connect trading

CREATE TABLE account_configs (
  "username" character varying(64) NOT NULL,
  "balance" numeric(19,3) NOT NULL,
  "limit_config" character varying(1024),
  CONSTRAINT "account_configs_username" UNIQUE ("username")
);

CREATE TABLE positions (
  "id" serial NOT NULL,
  "username" character varying(64) NOT NULL,
  "symbol" character varying(16) NOT NULL,
  "price" numeric(19,3) NOT NULL,
  "quantity" numeric(19,3) NOT NULL,
  "stock_price" numeric(19,3) NOT NULL,
  "updated_at" timestamp NOT NULL
);
