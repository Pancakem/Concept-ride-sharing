CREATE TABLE request (  
    id VARCHAR NOT NULL PRIMARY KEY,
    driverid VARCHAR NOT NULL ,
    riderid VARCHAR NOT NULL,
    origin VARCHAR,
    destination VARCHAR,
    actual_price REAL,
    completed BOOLEAN NOT NULL,
    ratings REAL,
    average_time DATETIME,
    create_date DATETIME default now,
);
