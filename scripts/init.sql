CREATE TABLE "waypoints" (
	"id" INTEGER NOT NULL UNIQUE,
	"symbol" VARCHAR NOT NULL UNIQUE,
	"system_id" VARCHAR,
	PRIMARY KEY("id")
);

CREATE INDEX "waypoints_index_0"
ON "waypoints" ("system_id");

CREATE UNIQUE INDEX "waypoints_index_1"
ON "waypoints" ("symbol");

CREATE TABLE "waypoints_traits" (
	"waypoint_id" INTEGER NOT NULL UNIQUE,
	"trait_id" INTEGER NOT NULL,
	PRIMARY KEY("waypoint_id", "trait_id")
);

CREATE UNIQUE INDEX "waypoints_traits_index_0"
ON "waypoints_traits" ("waypoint_id", "trait_id");

CREATE TABLE "traits" (
	"id" INTEGER NOT NULL UNIQUE,
	"symbol" VARCHAR NOT NULL UNIQUE,
	PRIMARY KEY("id")
);


CREATE TABLE "products" (
	"id" INTEGER NOT NULL UNIQUE,
	"symbol" VARCHAR NOT NULL UNIQUE,
	PRIMARY KEY("id")
);


CREATE TABLE "waypoints_products" (
	"waypoint_id" INTEGER NOT NULL UNIQUE,
	"product_id" INTEGER NOT NULL,
	"exports" BOOLEAN,
	"exchange" BOOLEAN,
	"imports" BOOLEAN,
	PRIMARY KEY("waypoint_id", "product_id")
);


ALTER TABLE "waypoints_traits"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_traits"
ADD FOREIGN KEY("trait_id") REFERENCES "traits"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_products"
ADD FOREIGN KEY("product_id") REFERENCES "products"("id")
ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE "waypoints_products"
ADD FOREIGN KEY("waypoint_id") REFERENCES "waypoints"("id")
ON UPDATE CASCADE ON DELETE CASCADE;