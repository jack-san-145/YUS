CREATE TABLE IF NOT EXISTS drivers (
    driver_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    mobile_no TEXT NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT ''
);
