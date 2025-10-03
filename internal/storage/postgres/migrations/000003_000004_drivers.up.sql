CREATE TABLE IF NOT EXISTS drivers (
    driver_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    mobile_no TEXT UNIQUE NOT NULL DEFAULT '',
    email TEXT UNIQUE NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT ''
);

-- Change the sequence to start from 1001
ALTER SEQUENCE drivers_driver_id_seq RESTART WITH 1001;
