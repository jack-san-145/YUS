
CREATE TABLE route_stops (
    route_id int NOT NULL,
    route_name TEXT NOT NULL DEFAULT '',
    direction text NOT NULL,
    stop_name text NOT NULL DEFAULT '',
    is_stop boolean NOT NULL DEFAULT false,
    lat text NOT NULL DEFAULT '',
    lon text NOT NULL DEFAULT '',
    arrival_time text NOT NULL DEFAULT '',
    departure_time text NOT NULL DEFAULT '',
    FOREIGN KEY (route_id, direction) REFERENCES all_routes(route_id, direction)
);

-- Create an index on (route_id, stop_name) for faster lookups
CREATE INDEX idx_routeid_stopname ON route_stops(route_id, stop_name,direction);
