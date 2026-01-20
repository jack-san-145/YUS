CREATE TABLE IF NOT EXISTS backup_routes (
    route_id int,
    route_name text,
    Path text,
    src text,
    dest text,
    direction text,
    route_json jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE backup_routes ADD PRIMARY KEY (route_id,route_name,src,dest,direction);