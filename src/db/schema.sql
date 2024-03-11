DROP DATABASE IF EXISTS db_platform;

CREATE DATABASE db_platform;

USE db_platform;

CREATE TABLE users(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    supervisor VARCHAR(255) DEFAULT "unknown",
    sector VARCHAR(255) DEFAULT "unknown",
    isAdmin BOOLEAN NOT NULL DEFAULT FALSE,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id, username),
    UNIQUE KEY (username) -- This ensures unique usernames
);

CREATE TABLE external_databases(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    sslMode VARCHAR(50) NOT NULL,
    userId BINARY(16) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name, host, id),
    -- This ensures unique combination of name, host, and id
    FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE TABLE logs(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    dbId BINARY(16),
    newUser VARCHAR(255) NOT NULL,
    wo INT NOT NULL,
    userId BINARY(16) NOT NULL,
    action VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (dbId) REFERENCES external_databases(id)
);

-- Insert the user and get its ID 
-- TODO change this password as the auth has changed
INSERT INTO
    users (username, password, isAdmin)
VALUES
    (
        "test",
        "JFVMdGtBaXgxcHVmdTlYeTFuV0hkckEkYjFpdUtJRHc2Z0o5cCtMeFh3THA5Yll4QitSVnNjaGJpK3VnY0paaGRyaw",
        TRUE
    );

-- Get the ID of the newly inserted user
SET
    @userId = (
        SELECT
            id
        FROM
            users
        WHERE
            username = "test"
    );

-- Insert into external_databases table using the obtained user ID
INSERT INTO
    external_databases (name, host, port, type, sslMode, userId)
VALUES
    (
        "mysql",
        "localhost",
        3001,
        "mysql",
        "disable",
        @userId
    ),
    (
        "mysql-2",
        "localhost",
        3003,
        "mysql",
        "disable",
        @userId
    ),
    (
        "mysql-3",
        "localhost",
        3004,
        "mysql",
        "disable",
        @userId
    ),
    (
        "maria",
        "localhost",
        3002,
        "mysql",
        "disable",
        @userId
    ),
    (
        "maria-2",
        "localhost",
        3005,
        "mysql",
        "disable",
        @userId
    ),
    (
        "postgres",
        "localhost",
        5432,
        "postgres",
        "disable",
        @userId
    ),
    (
        "postgres-2",
        "localhost",
        5433,
        "postgres",
        "disable",
        @userId
    );