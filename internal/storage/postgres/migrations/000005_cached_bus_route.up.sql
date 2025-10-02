CREATE TABLE IF NOT EXISTS cached_bus_route (
    bus_id int  NOT NULL REFERENCES current_bus_route(bus_id),
    driver_id int  NOT NULL REFERENCES drivers(driver_id),
    route_id int NOT NULL 
);
