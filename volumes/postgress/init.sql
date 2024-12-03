CREATE TABLE IF NOT EXISTS my_table (
    unique_id VARCHAR(255) PRIMARY KEY,   -- unique_id as a string with a maximum length of 255 characters
    timestamp TIMESTAMP NOT NULL          -- timestamp column, no default value
);
