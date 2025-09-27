CREATE TABLE IF NOT EXISTS all_routes (
    route_id SERIAL PRIMARY KEY,
    src TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT '',
    arrival_time TEXT NOT NULL DEFAULT '',
    departure_time TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--creating index for source and destination
CREATE INDEX idx_src_dest ON all_routes(src, dest);
