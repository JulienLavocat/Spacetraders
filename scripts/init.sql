CREATE TABLE "waypoints" (
	"id" VARCHAR NOT NULL,
	"system_id" VARCHAR NOT NULL,
	"x" INTEGER NOT NULL,
	"y" INTEGER NOT NULL,
	"type" VARCHAR NOT NULL,
	"faction" VARCHAR NOT NULL,
	"orbits" VARCHAR,
	"under_construction" BOOLEAN NOT NULL,
	"submitted_on" TIMESTAMPTZ,
	"submitted_by" VARCHAR,
	"gid" SERIAL NOT NULL UNIQUE DEFAULT nextval('waypoints_gid_seq'),
	PRIMARY KEY("id")
);

CREATE INDEX "waypoints_index_0"
ON "waypoints" ("system_id");

CREATE TABLE "waypoints_traits" (
	"waypoint_id" VARCHAR NOT NULL,
	"trait_id" VARCHAR NOT NULL,
	PRIMARY KEY("waypoint_id", "trait_id")
);


CREATE TABLE "traits" (
	"id" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE "products" (
	"id" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE "waypoints_products" (
	"waypoint_id" VARCHAR NOT NULL,
	"product_id" VARCHAR NOT NULL,
	"export" BOOLEAN NOT NULL DEFAULT false,
	"exchange" BOOLEAN NOT NULL DEFAULT false,
	"import" BOOLEAN NOT NULL DEFAULT false,
	"volume" INTEGER DEFAULT null,
	"supply" VARCHAR DEFAULT null,
	"activity" VARCHAR DEFAULT null,
	"buy" INTEGER DEFAULT null,
	"sell" INTEGER DEFAULT null,
	"updated_at" TIMESTAMPTZ DEFAULT null,
	PRIMARY KEY("waypoint_id", "product_id")
);

CREATE UNIQUE INDEX "waypoints_products_index_0"
ON "waypoints_products" ("waypoint_id", "product_id", "export", "exchange", "import");

CREATE TABLE "systems" (
	"id" VARCHAR NOT NULL,
	"sector_id" VARCHAR NOT NULL,
	"type" VARCHAR NOT NULL,
	"x" INTEGER NOT NULL,
	"y" INTEGER NOT NULL,
	"gid" SERIAL NOT NULL UNIQUE DEFAULT nextval('waypoints_gid_seq'),
	PRIMARY KEY("id")
);


CREATE TABLE "factions" (
	"id" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE "factions_systems" (
	"faction_id" VARCHAR NOT NULL,
	"system_id" VARCHAR NOT NULL,
	PRIMARY KEY("faction_id", "system_id")
);


CREATE TABLE "factions_waypoints" (
	"waypoint_id" VARCHAR NOT NULL,
	"faction_id" VARCHAR NOT NULL,
	PRIMARY KEY("waypoint_id", "faction_id")
);


CREATE TABLE "modifiers" (
	"id" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE "waypoints_modifiers" (
	"modifier_id" VARCHAR NOT NULL,
	"waypoint_id" VARCHAR NOT NULL,
	PRIMARY KEY("modifier_id", "waypoint_id")
);


CREATE TABLE "market_probes" (
	"ship" VARCHAR,
	"waypoint" VARCHAR,
	"system" VARCHAR NOT NULL,
	PRIMARY KEY("ship", "waypoint")
);


CREATE TABLE "waypoints_graphs" (
	"id" SERIAL NOT NULL UNIQUE,
	"source" SERIAL NOT NULL,
	"target" SERIAL NOT NULL,
	"cost" INTEGER NOT NULL,
	"system_id" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);

CREATE INDEX "waypoints_graphs_index_0"
ON "waypoints_graphs" ("system_id");

CREATE INDEX "waypoints_graphs_index_1"
ON "waypoints_graphs" ("cost");

CREATE INDEX "waypoints_graphs_index_2"
ON "waypoints_graphs" ("source");

CREATE INDEX "waypoints_graphs_index_3"
ON "waypoints_graphs" ("target");

ALTER TABLE "factions_systems"
ADD FOREIGN KEY("faction_id") REFERENCES "factions"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "factions_systems"
ADD FOREIGN KEY("system_id") REFERENCES "systems"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_modifiers"
ADD FOREIGN KEY("modifier_id") REFERENCES "modifiers"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_modifiers"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "factions_waypoints"
ADD FOREIGN KEY("faction_id") REFERENCES "factions"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "factions_waypoints"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_products"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_products"
ADD FOREIGN KEY("product_id") REFERENCES "products"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_traits"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_traits"
ADD FOREIGN KEY("trait_id") REFERENCES "traits"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "waypoints"
ADD FOREIGN KEY("system_id") REFERENCES "systems"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_graphs"
ADD FOREIGN KEY("system_id") REFERENCES "systems"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "waypoints_graphs"
ADD FOREIGN KEY("source") REFERENCES "waypoints"("gid")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "waypoints_graphs"
ADD FOREIGN KEY("target") REFERENCES "waypoints"("gid")
ON UPDATE NO ACTION ON DELETE NO ACTION;