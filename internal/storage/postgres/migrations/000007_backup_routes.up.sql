CREATE TABLE IF NOT EXISTS backup_routes (
    route_id int,
    src text,
    dest text,
    route_name text,
    direction text,
    route_json jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE backup_routes ADD PRIMARY KEY (route_id, src,dest,route_name,direction);