CREATE TABLE route_stops (
    route_id int REFERENCES all_routes(route_id),
    stop_name text NOT NULL DEFAULT '',
    is_stop boolean NOT NULL DEFAULT false,
    arrival_time text NOT NULL DEFAULT '',
    departure_time text NOT NULL DEFAULT ''
);

-- Create an index on route_id and stop_name
CREATE INDEX idx_routeid_stopname ON route_stops(route_id, stop_name);
