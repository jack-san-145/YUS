CREATE TABLE driver_removal_request (
    driver_id INT PRIMARY KEY,
    driver_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
