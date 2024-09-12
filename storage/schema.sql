CREATE USER 'apt_db_platform'@'%' IDENTIFIED BY '1qaz!EDC';

GRANT ALL PRIVILEGES ON *.* TO 'apt_db_platform'@'%' WITH GRANT OPTION;

FLUSH PRIVILEGES;

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
    PRIMARY KEY (id),
    UNIQUE KEY (username) -- This ensures unique usernames
);

CREATE TABLE external_databases(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    sslMode VARCHAR(50) NOT NULL,
    createdBy BINARY(16) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name, host, id), -- This ensures unique combination of name, host, and id
    UNIQUE KEY (name),
    FOREIGN KEY (createdBy) REFERENCES users(id)
);

CREATE TABLE logs(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    dbId BINARY(16),
    newUser VARCHAR(255) NOT NULL,
    wo INT NOT NULL,
    createdBy BINARY(16) NOT NULL,
    action VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (createdBy) REFERENCES users(id),
    FOREIGN KEY (dbId) REFERENCES external_databases(id)
);

-- Insert the default admin user with password test
INSERT INTO
    users (username, password, isAdmin)
VALUES
    (
        "admin",
        "JFVMdGtBaXgxcHVmdTlYeTFuV0hkckEkYjFpdUtJRHc2Z0o5cCtMeFh3THA5Yll4QitSVnNjaGJpK3VnY0paaGRyaw",
        TRUE
    );

-- Add username and password fields to external_databases table
ALTER TABLE external_databases
ADD COLUMN username VARCHAR(255) NOT NULL,
ADD COLUMN password VARCHAR(255) NOT NULL;

-- Create administrational logs table
CREATE TABLE admin_logs(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    action VARCHAR(255) NOT NULL,
    resourceId BINARY(16) NOT NULL,
    resourceType VARCHAR(255) NOT NULL,
    userId BINARY(16) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY (userId) REFERENCES users(id)
);

-- Add owner to external_databases
ALTER TABLE external_databases
ADD COLUMN owner VARCHAR(255) NOT NULL DEFAULT "verificar";

-- Make dbId not null on logs table
ALTER TABLE logs
MODIFY COLUMN dbId BINARY(16) NOT NULL;
