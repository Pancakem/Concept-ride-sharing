CREATE TABLE request (  
    id VARCHAR NOT NULL PRIMARY KEY,
    driverid VARCHAR NOT NULL ,
    riderid VARCHAR NOT NULL,
    actual_price REAL,
    completed BOOLEAN NOT NULL,
    ride_time VARCHAR,
    distance REAL,
    -- how long it takes to match
    average_time DATETIME, 
    create_date DATETIME default now,
);

CREATE TABLE place (
    id VARCHAR NOT NULL,
    -- ORIGIN
    origin_name VARCHAR,
    origin_latitude REAL,
    origin_longitude REAL,
    -- DESTINATION
    destination_name VARCHAR,
    destination_latitude REAL,
    destination_longitude REAL,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES request(id),
); 

CREATE TABLE rating (
    rideid VARCHAR NOT NULL,
    riderrating REAL,
    driverrating REAL,
);