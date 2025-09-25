CREATE TABLE bus_route (
    bus_id int PRIMARY KEY,
    route_id int REFERENCES all_routes(route_id)
);
