CREATE TABLE request (  
    id VARCHAR (50) NOT NULL PRIMARY KEY,
    driverid VARCHAR (50) NOT NULL ,
    riderid VARCHAR (50) NOT NULL,
    actual_price REAL,
    completed BOOLEAN NOT NULL,
    ride_time VARCHAR (50),
    distance REAL,
    average_time VARCHAR (50), 
    create_date VARCHAR (50)
);

CREATE TABLE place (
    id VARCHAR NOT NULL PRIMARY KEY ,
    origin_name VARCHAR (50),
    origin_latitude REAL,
    origin_longitude REAL,
    destination_name VARCHAR (50),
    destination_latitude REAL,
    destination_longitude REAL,
    
    FOREIGN KEY (id) REFERENCES request (id)
); 

CREATE TABLE ratings (
    rideid VARCHAR (50) NOT NULL,
    riderrating REAL,
    driverrating REAL
);