CREATE TABLE all_routes (
    route_id SERIAL PRIMARY KEY,
    src TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    departure_time TEXT NOT NULL DEFAULT '',
    arrival_time TEXT NOT NULL DEFAULT ''
);

-- Composite index on (src, dest)
CREATE INDEX idx_src_dest ON all_routes(src, dest);
