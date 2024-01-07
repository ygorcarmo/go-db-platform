DROP DATABASE IF EXISTS db_platform;

CREATE DATABASE db_platform;

USE db_platform;

CREATE TABLE db_connection_info (
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    sslMode VARCHAR(50) NOT NULL,
    PRIMARY KEY (name, host)
);
