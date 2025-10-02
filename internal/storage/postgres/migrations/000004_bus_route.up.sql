CREATE TABLE IF NOT EXISTS bus_route (
    bus_id int PRIMARY KEY,
    driver_id int REFERENCES drivers(driver_id),
    route_id int NOT NULL 
);
