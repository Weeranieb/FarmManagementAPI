CREATE TABLE "Clients" (
  "Id" serial PRIMARY KEY,
  "Name" varchar NOT NULL,
  "OwnerName" varchar NOT NULL,
  "ContactNumber" varchar NOT NULL,
  "IsActive" boolean NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Users" (
  "Id" serial PRIMARY KEY,
  "ClientId" bigint NOT NULL,
  "Username" varchar NOT NULL,
  "Password" varchar NOT NULL,
  "FirstName" varchar NOT NULL,
  "LastName" varchar,
  "ContactNumber" varchar NOT NULL,
  "IsAdmin" boolean NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Farms" (
  "Id" serial PRIMARY KEY,
  "ClientId" bigint NOT NULL,
  "Code" varchar NOT NULL,
  "Name" varchar NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "FarmGroups" (
  "Id" serial PRIMARY KEY,
  "ClientId" bigint NOT NULL,
  "Code" varchar NOT NULL,
  "Name" varchar NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "FarmOnFarmGroup" (
  "Id" serial PRIMARY KEY,
  "FarmId" bigint NOT NULL,
  "FarmGroupId" bigint NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Ponds" (
  "Id" serial PRIMARY KEY,
  "FarmId" bigint NOT NULL,
  "Code" varchar NOT NULL,
  "Name" varchar NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "ActivePonds" (
  "Id" serial PRIMARY KEY,
  "PondId" bigint NOT NULL,
  "StartDate" date NOT NULL,
  "EndDate" date,
  "IsActive" boolean NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Activities" (
  "Id" serial PRIMARY KEY,
  "ActivePondId" bigint NOT NULL,
  "ToActivePondId" bigint,
  "Mode" varchar NOT NULL,
  "MerchantId" bigint,
  "Amount" integer,
  "FishType" varchar,
  "FishWeight" float,
  "PricePerUnit" float,
  "ActivityDate" date NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "AdditionalCosts" (
  "Id" serial PRIMARY KEY,
  "ActivityId" bigint NOT NULL,
  "Title" varchar NOT NULL,
  "Cost" float NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Merchants" (
  "Id" serial PRIMARY KEY,
  "Name" varchar NOT NULL,
  "ContactNumber" varchar NOT NULL,
  "Location" varchar NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "SellDetails" (
  "Id" serial PRIMARY KEY,
  "SellId" bigint NOT NULL,
  "Size" varchar NOT NULL,
  "FishType" varchar,
  "Amount" float NOT NULL,
  "FishUnit" varchar NOT NULL,
  "PricePerUnit" float NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Bills" (
  "Id" serial PRIMARY KEY,
  "Type" varchar NOT NULL,
  "Other" varchar,
  "FarmGroupId" integer NOT NULL,
  "PaidAmount" float NOT NULL,
  "PaymentDate" date NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "Workers" (
  "Id" serial PRIMARY KEY,
  "ClientId" bigint NOT NULL,
  "FarmGroupId" bigint NOT NULL,
  "FirstName" varchar NOT NULL,
  "LastName" varchar,
  "ContactNumber" varchar,
  "Salary" integer NOT NULL,
  "HireDate" date,
  "IsActive" boolean NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "FeedCollections" (
  "Id" serial PRIMARY KEY,
  "ClientId" bigint NOT NULL,
  "Code" varchar NOT NULL,
  "Name" varchar NOT NULL,
  "Unit" varchar NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "DailyFeeds" (
  "Id" serial PRIMARY KEY,
  "ActivePondId" bigint NOT NULL,
  "FeedCollectionId" bigint NOT NULL,
  "Amount" float NOT NULL,
  "FeedDate" date NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE TABLE "FeedPriceHistories" (
  "Id" serial PRIMARY KEY,
  "FeedCollectionId" bigint NOT NULL,
  "Price" float NOT NULL,
  "PriceUpdatedDate" date NOT NULL,
  "DelFlag" boolean NOT NULL,
  "CreatedDate" timestamp NOT NULL DEFAULT (now()),
  "CreatedBy" varchar NOT NULL,
  "UpdatedDate" timestamp NOT NULL DEFAULT (now()),
  "UpdatedBy" varchar NOT NULL
);

CREATE INDEX ON "Users" ("ClientId");

CREATE INDEX ON "Farms" ("ClientId");

CREATE INDEX ON "FarmGroups" ("ClientId");

CREATE INDEX ON "FarmOnFarmGroup" ("FarmGroupId");

CREATE INDEX ON "FarmOnFarmGroup" ("FarmId");

CREATE INDEX ON "Ponds" ("FarmId");

CREATE INDEX ON "ActivePonds" ("PondId");

CREATE INDEX ON "Activities" ("ActivePondId");

CREATE INDEX ON "Activities" ("MerchantId");

CREATE INDEX ON "AdditionalCosts" ("ActivityId");

CREATE INDEX ON "Bills" ("FarmGroupId");

CREATE INDEX ON "Workers" ("ClientId");

CREATE INDEX ON "Workers" ("FarmGroupId");

CREATE INDEX ON "FeedCollections" ("ClientId");

CREATE INDEX ON "DailyFeeds" ("ActivePondId");

CREATE INDEX ON "DailyFeeds" ("FeedCollectionId");

CREATE INDEX ON "FeedPriceHistories" ("FeedCollectionId");

ALTER TABLE "Users" ADD FOREIGN KEY ("ClientId") REFERENCES "Clients" ("Id");

ALTER TABLE "Farms" ADD FOREIGN KEY ("ClientId") REFERENCES "Clients" ("Id");

ALTER TABLE "FarmGroups" ADD FOREIGN KEY ("ClientId") REFERENCES "Clients" ("Id");

ALTER TABLE "FarmOnFarmGroup" ADD FOREIGN KEY ("FarmId") REFERENCES "Farms" ("Id");

ALTER TABLE "FarmOnFarmGroup" ADD FOREIGN KEY ("FarmGroupId") REFERENCES "FarmGroups" ("Id");

ALTER TABLE "Ponds" ADD FOREIGN KEY ("FarmId") REFERENCES "Farms" ("Id");

ALTER TABLE "ActivePonds" ADD FOREIGN KEY ("PondId") REFERENCES "Ponds" ("Id");

ALTER TABLE "Activities" ADD FOREIGN KEY ("ActivePondId") REFERENCES "ActivePonds" ("Id");

ALTER TABLE "Activities" ADD FOREIGN KEY ("ToActivePondId") REFERENCES "ActivePonds" ("Id");

ALTER TABLE "Activities" ADD FOREIGN KEY ("MerchantId") REFERENCES "Merchants" ("Id");

ALTER TABLE "AdditionalCosts" ADD FOREIGN KEY ("ActivityId") REFERENCES "Activities" ("Id");

ALTER TABLE "SellDetails" ADD FOREIGN KEY ("SellId") REFERENCES "Activities" ("Id");

ALTER TABLE "Bills" ADD FOREIGN KEY ("FarmGroupId") REFERENCES "FarmGroups" ("Id");

ALTER TABLE "Workers" ADD FOREIGN KEY ("ClientId") REFERENCES "Clients" ("Id");

ALTER TABLE "Workers" ADD FOREIGN KEY ("FarmGroupId") REFERENCES "FarmGroups" ("Id");

ALTER TABLE "FeedCollections" ADD FOREIGN KEY ("ClientId") REFERENCES "Clients" ("Id");

ALTER TABLE "DailyFeeds" ADD FOREIGN KEY ("ActivePondId") REFERENCES "ActivePonds" ("Id");

ALTER TABLE "DailyFeeds" ADD FOREIGN KEY ("FeedCollectionId") REFERENCES "FeedCollections" ("Id");

ALTER TABLE "FeedPriceHistories" ADD FOREIGN KEY ("FeedCollectionId") REFERENCES "FeedCollections" ("Id");
