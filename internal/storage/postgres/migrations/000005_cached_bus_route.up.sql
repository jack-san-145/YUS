CREATE TABLE IF NOT EXISTS cached_bus_route (
    bus_id int REFERENCES current_bus_route(bus_id),
    driver_id int REFERENCES drivers(driver_id),
    route_id int NOT NULL 
);
