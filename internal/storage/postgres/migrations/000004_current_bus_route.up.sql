CREATE TABLE IF NOT EXISTS current_bus_route (
    bus_id int PRIMARY KEY,
    driver_id int NOT NULL,
    route_id int NOT NULL 
);
