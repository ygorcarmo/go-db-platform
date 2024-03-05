DROP DATABASE IF EXISTS db_platform;

CREATE DATABASE db_platform;

USE db_platform;

CREATE TABLE users(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    isAdmin BOOLEAN NOT NULL DEFAULT FALSE, 
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id, username)
);

CREATE TABLE db_connection_info (
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    sslMode VARCHAR(50) NOT NULL,
    userId BINARY(16),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (name, host, id),
    FOREIGN KEY (userId) REFERENCES users(id)
);

INSERT INTO
    users (username, password, isAdmin)
VALUES
    (
        "test",
        "JGFyZ29uMmlkJHY9MTkkbT02NTUzNix0PTMscD0yJFovVGJ0Q3V4b0dBOWowQzNsR2ttK0EkaGhqVjJGdTQ3WWdFU2RqSm1BQTN6ZTNsM2ZtaGQvcVduNWlscVQ2THE4OA",
        TRUE
    );

INSERT INTO
    db_connection_info (name, host, port, type, sslMode)
VALUES
    ("mysql", "localhost", 3001, "mysql", "disable"),
    ("mysql-2", "localhost", 3003, "mysql", "disable"),
    ("mysql-3", "localhost", 3004, "mysql", "disable"),
    ("maria", "localhost", 3002, "mysql", "disable"),
    ("maria-2", "localhost", 3005, "mysql", "disable"),
    (
        "postgres",
        "localhost",
        5432,
        "postgres",
        "disable"
    ),
    (
        "postgres-2",
        "localhost",
        5433,
        "postgres",
        "disable"
    );