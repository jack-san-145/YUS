
CREATE TABLE IF NOT EXISTS all_routes (
    route_id int NOT NULL,
    route_name TEXT NOT NULL DEFAULT '',
    src TEXT NOT NULL DEFAULT '',
    dest TEXT NOT NULL DEFAULT '',
    direction text NOT NULL DEFAULT '',
    departure_time TEXT NOT NULL DEFAULT '',
    arrival_time TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (route_id, direction)
);

CREATE INDEX idx_src_dest ON all_routes(src, dest);

ALTER TABLE all_routes
ADD CONSTRAINT unique_name_src_dest_direction 
UNIQUE (route_name, src, dest, direction);
