CREATE TABLE IF NOT EXISTS current_bus_route (
    bus_id int PRIMARY KEY,
    driver_id int NOT NULL DEFAULT 1000,
    route_id int NOT NULL DEFAULT 0,
    direction TEXT NOT NULL DEFAULT 'UP',
    route_name TEXT NOT NULL DEFAULT '',
    src TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT ''
);
