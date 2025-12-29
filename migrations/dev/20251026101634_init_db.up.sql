CREATE TABLE clients (
  id bigint PRIMARY KEY,
  name varchar NOT NULL,
  owner_name varchar NOT NULL,
  contact_number varchar NOT NULL,
  is_active bit NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE users (
  id bigint PRIMARY KEY,
  client_id bigint NOT NULL,
  username varchar NOT NULL,
  password varchar NOT NULL,
  first_name varchar NOT NULL,
  last_name varchar,
  contact_number varchar NOT NULL,
  user_level integer NOT NULL DEFAULT 0,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE farms (
  id bigint PRIMARY KEY,
  client_id bigint NOT NULL,
  code varchar NOT NULL,
  name varchar NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE farm_groups (
  id bigint PRIMARY KEY,
  client_id bigint NOT NULL,
  code varchar NOT NULL,
  name varchar NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE farm_on_farm_group (
  id bigint PRIMARY KEY,
  farm_id bigint NOT NULL,
  farm_group_id bigint NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE ponds (
  id bigint PRIMARY KEY,
  farm_id bigint NOT NULL,
  code varchar NOT NULL,
  name varchar NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE active_ponds (
  id bigint PRIMARY KEY,
  pond_id bigint NOT NULL,
  start_date date NOT NULL,
  end_date date,
  is_active bit NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE activities (
  id bigint PRIMARY KEY,
  active_pond_id bigint NOT NULL,
  to_active_pond_id bigint,
  mode varchar NOT NULL,
  merchant_id bigint,
  amount integer,
  fish_type varchar,
  fish_weight float,
  fish_unit varchar,
  price_per_unit float,
  activity_date date NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE additional_costs (
  id bigint PRIMARY KEY,
  activity_id bigint NOT NULL,
  title varchar NOT NULL,
  cost float NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE merchants (
  id bigint PRIMARY KEY,
  name varchar NOT NULL,
  contact_number varchar NOT NULL,
  location varchar NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE sell_details (
  id bigint PRIMARY KEY,
  sell_id bigint NOT NULL,
  size varchar NOT NULL,
  fish_type varchar,
  amount float NOT NULL,
  fish_unit varchar NOT NULL,
  price_per_unit float NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE bills (
  id bigint PRIMARY KEY,
  type varchar NOT NULL,
  other varchar,
  farm_group_id integer NOT NULL,
  paid_amount float NOT NULL,
  payment_date date NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE workers (
  id bigint PRIMARY KEY,
  client_id bigint NOT NULL,
  farm_group_id bigint NOT NULL,
  first_name varchar NOT NULL,
  last_name varchar,
  contact_number varchar,
  nationality varchar NOT NULL DEFAULT '',
  salary integer NOT NULL,
  hire_date date,
  is_active bit NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE feed_collections (
  id bigint PRIMARY KEY,
  client_id bigint NOT NULL,
  code varchar NOT NULL,
  name varchar NOT NULL,
  unit varchar NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE daily_feeds (
  id bigint PRIMARY KEY,
  active_pond_id bigint,
  pond_id bigint NOT NULL,
  feed_collection_id bigint NOT NULL,
  amount float NOT NULL,
  feed_date date NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE TABLE feed_price_histories (
  id bigint PRIMARY KEY,
  feed_collection_id bigint NOT NULL,
  price float NOT NULL,
  price_updated_date date NOT NULL,
  deleted_at timestamp,
  created_at timestamp NOT NULL DEFAULT (now()),
  created_by varchar NOT NULL,
  updated_at timestamp NOT NULL DEFAULT (now()),
  updated_by varchar NOT NULL
);

CREATE INDEX ON users (client_id);

CREATE INDEX ON farms (client_id);

CREATE INDEX ON farm_groups (client_id);

CREATE INDEX ON farm_on_farm_group (farm_group_id);

CREATE INDEX ON farm_on_farm_group (farm_id);

CREATE INDEX ON ponds (farm_id);

CREATE INDEX ON active_ponds (pond_id);

CREATE INDEX ON activities (active_pond_id);

CREATE INDEX ON activities (merchant_id);

CREATE INDEX ON additional_costs (activity_id);

CREATE INDEX ON bills (farm_group_id);

CREATE INDEX ON workers (client_id);

CREATE INDEX ON workers (farm_group_id);

CREATE INDEX ON feed_collections (client_id);

CREATE INDEX ON daily_feeds (active_pond_id);

CREATE INDEX ON daily_feeds (feed_collection_id);

CREATE INDEX ON daily_feeds (pond_id);

CREATE INDEX ON feed_price_histories (feed_collection_id);

ALTER TABLE users ADD FOREIGN KEY (client_id) REFERENCES clients (id);

ALTER TABLE farms ADD FOREIGN KEY (client_id) REFERENCES clients (id);

ALTER TABLE farm_groups ADD FOREIGN KEY (client_id) REFERENCES clients (id);

ALTER TABLE farm_on_farm_group ADD FOREIGN KEY (farm_id) REFERENCES farms (id);

ALTER TABLE farm_on_farm_group ADD FOREIGN KEY (farm_group_id) REFERENCES farm_groups (id);

ALTER TABLE ponds ADD FOREIGN KEY (farm_id) REFERENCES farms (id);

ALTER TABLE active_ponds ADD FOREIGN KEY (pond_id) REFERENCES ponds (id);

ALTER TABLE activities ADD FOREIGN KEY (active_pond_id) REFERENCES active_ponds (id);

ALTER TABLE activities ADD FOREIGN KEY (to_active_pond_id) REFERENCES active_ponds (id);

ALTER TABLE activities ADD FOREIGN KEY (merchant_id) REFERENCES merchants (id);

ALTER TABLE additional_costs ADD FOREIGN KEY (activity_id) REFERENCES activities (id);

ALTER TABLE sell_details ADD FOREIGN KEY (sell_id) REFERENCES activities (id);

ALTER TABLE bills ADD FOREIGN KEY (farm_group_id) REFERENCES farm_groups (id);

ALTER TABLE workers ADD FOREIGN KEY (client_id) REFERENCES clients (id);

ALTER TABLE workers ADD FOREIGN KEY (farm_group_id) REFERENCES farm_groups (id);

ALTER TABLE feed_collections ADD FOREIGN KEY (client_id) REFERENCES clients (id);

ALTER TABLE daily_feeds ADD FOREIGN KEY (active_pond_id) REFERENCES active_ponds (id);

ALTER TABLE daily_feeds ADD FOREIGN KEY (pond_id) REFERENCES ponds (id);

ALTER TABLE daily_feeds ADD FOREIGN KEY (feed_collection_id) REFERENCES feed_collections (id);

ALTER TABLE feed_price_histories ADD FOREIGN KEY (feed_collection_id) REFERENCES feed_collections (id);
