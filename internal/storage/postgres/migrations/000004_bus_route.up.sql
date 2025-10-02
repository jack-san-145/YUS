CREATE TABLE IF NOT EXISTS current_bus_route (
    bus_id int PRIMARY KEY,
    driver_id int REFERENCES drivers(driver_id),
    route_id int NOT NULL 
);
