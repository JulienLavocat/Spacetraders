
create table if not exists public.factions
(
    id varchar not null
        primary key
);

alter table public.factions
    owner to spacetraders;

create table if not exists public.market_probes
(
    ship     varchar               not null,
    waypoint varchar               not null,
    system   varchar               not null,
    shipyard boolean default false not null,
    constraint market_probes_pk
        primary key (ship, waypoint)
);

alter table public.market_probes
    owner to spacetraders;

create table if not exists public.modifiers
(
    id varchar not null
        primary key
);

alter table public.modifiers
    owner to spacetraders;

create table if not exists public.products
(
    id varchar not null
        primary key
);

alter table public.products
    owner to spacetraders;

create table if not exists public.systems
(
    id        varchar                                                                 not null
        primary key,
    sector_id varchar                                                                 not null,
    type      varchar                                                                 not null,
    x         integer                                                                 not null,
    y         integer                                                                 not null,
    gid       serial,
    geom      geometry default '010100000000000000000000000000000000000000'::geometry not null
);

alter table public.systems
    owner to spacetraders;

create table if not exists public.factions_systems
(
    faction_id varchar not null
        references public.factions
            on update cascade on delete cascade,
    system_id  varchar not null
        references public.systems
            on update cascade on delete cascade,
    primary key (faction_id, system_id)
);

alter table public.factions_systems
    owner to spacetraders;

create table if not exists public.traits
(
    id varchar not null
        primary key
);

alter table public.traits
    owner to spacetraders;

create table if not exists public.waypoints
(
    id                 varchar                                                                 not null
        primary key,
    system_id          varchar                                                                 not null
        references public.systems
            on update cascade on delete cascade,
    x                  integer                                                                 not null,
    y                  integer                                                                 not null,
    type               varchar                                                                 not null,
    faction            varchar                                                                 not null,
    orbits             varchar,
    under_construction boolean                                                                 not null,
    submitted_on       timestamp with time zone,
    submitted_by       varchar,
    geom               geometry default '010100000000000000000000000000000000000000'::geometry not null,
    gid                integer  default nextval('waypoints_gid_seq1'::regclass)                not null
);

alter table public.waypoints
    owner to spacetraders;

create table if not exists public.factions_waypoints
(
    waypoint_id varchar not null
        references public.waypoints
            on update cascade on delete cascade,
    faction_id  varchar not null
        references public.factions
            on update cascade on delete cascade,
    primary key (waypoint_id, faction_id)
);

alter table public.factions_waypoints
    owner to spacetraders;

create index if not exists waypoints_index_0
    on public.waypoints (system_id);

create table if not exists public.waypoints_modifiers
(
    modifier_id varchar not null
        references public.modifiers
            on update cascade on delete cascade,
    waypoint_id varchar not null
        references public.waypoints
            on update cascade on delete cascade,
    primary key (modifier_id, waypoint_id)
);

alter table public.waypoints_modifiers
    owner to spacetraders;

create table if not exists public.waypoints_products
(
    waypoint_id varchar               not null
        references public.waypoints
            on update cascade on delete cascade,
    product_id  varchar               not null
        references public.products
            on update cascade on delete cascade,
    export      boolean default false not null,
    exchange    boolean default false not null,
    import      boolean default false not null,
    volume      integer,
    supply      varchar,
    activity    varchar,
    buy         integer,
    sell        integer,
    updated_at  timestamp with time zone,
    primary key (waypoint_id, product_id)
);

alter table public.waypoints_products
    owner to spacetraders;

create unique index if not exists waypoints_products_index_0
    on public.waypoints_products (waypoint_id, product_id, export, exchange, import);

create table if not exists public.waypoints_traits
(
    waypoint_id varchar not null
        references public.waypoints
            on update cascade on delete cascade,
    trait_id    varchar not null
        references public.traits,
    primary key (waypoint_id, trait_id)
);

alter table public.waypoints_traits
    owner to spacetraders;

create table if not exists public.waypoints_graphs
(
    id        serial
        primary key,
    source    serial,
    target    serial,
    cost      integer not null,
    system_id varchar not null
);

alter table public.waypoints_graphs
    owner to spacetraders;

create index if not exists waypoints_graphs_index_0
    on public.waypoints_graphs (system_id);

create index if not exists waypoints_graphs_index_1
    on public.waypoints_graphs (cost);

create index if not exists waypoints_graphs_index_2
    on public.waypoints_graphs (source);

create index if not exists waypoints_graphs_index_3
    on public.waypoints_graphs (target);

create table if not exists public.transactions
(
    id             serial
        constraint transactions_pk
            primary key,
    waypoint       varchar                                not null
        constraint transactions_waypoints_id_fk
            references public.waypoints,
    product        varchar                                not null
        constraint transactions_products_id_fk
            references public.products,
    amount         integer                                not null,
    type           varchar                                not null,
    ship           varchar                                not null,
    price_per_unit integer                                not null,
    total_price    integer                                not null,
    timestamp      timestamp with time zone default now() not null,
    agent_balance  bigint                                 not null
);

alter table public.transactions
    owner to spacetraders;

create index if not exists transactions_product_index
    on public.transactions (product);

create index if not exists transactions_ship_index
    on public.transactions (ship);

create index if not exists transactions_waypoint_index
    on public.transactions (waypoint);

