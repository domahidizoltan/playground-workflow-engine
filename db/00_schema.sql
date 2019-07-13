CREATE TABLE "stock_data" (
  "id" serial NOT NULL,
  "symbol" character varying(16) NOT NULL,
  "price" numeric(19,3) NOT NULL,
  "date" timestamp NOT NULL
);