CREATE TABLE IF NOT EXISTS cached_bus_route (
    bus_id int  NOT NULL REFERENCES current_bus_route(bus_id),
    route_id int NOT NULL DEFAULT 0,
    route_name TEXT NOT NULL DEFAULT '',
    src TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_cached_bus_route ON cached_bus_route(bus_id,route_id);