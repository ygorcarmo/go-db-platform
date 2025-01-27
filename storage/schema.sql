CREATE USER 'apt_db_platform' @'%' IDENTIFIED BY '1qaz!EDC';

GRANT ALL PRIVILEGES ON *.* TO 'apt_db_platform' @'%' WITH GRANT OPTION;

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
    loginAttempts INT NOT NULL DEFAULT 0,
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
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL DEFAULT "verificar",
    createdBy VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    protocol CHAR(4) DEFAULT "tcp",
    host_fallback VARCHAR(255) DEFAULT "localhost",
    port_fallback INT DEFAULT 1521,
    protocol_fallback CHAR(4) DEFAULT "tcp",
    PRIMARY KEY (id),
    UNIQUE KEY (name, host, id),
    -- This ensures unique combination of name, host, and id
    UNIQUE KEY (name)
);

CREATE TABLE logs(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    newUser VARCHAR(255) NOT NULL,
    wo INT NOT NULL,
    createdBy VARCHAR(255),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dbName VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    success BOOLEAN NOT NULL,
    PRIMARY KEY(id)
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

-- Create administrational logs table
CREATE TABLE admin_logs(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    action VARCHAR(255) NOT NULL,
    resourceType VARCHAR(255) NOT NULL,
    resourceName VARCHAR(255) NOT NULL DEFAULT "TBC",
    username VARCHAR(255) NOT NULL DEFAULT "TBC",
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

-- New feature: AD Authentication
-- TODO Drop DbId FK in logs, and add username to logs and admin logs as we wont add the ad users to DB
-- Add Config Table
CREATE TABLE ldap_config(
    id TINYINT UNSIGNED NOT NULL PRIMARY KEY DEFAULT 1,
    connectionStr VARCHAR(255) NOT NULL DEFAULT "127.0.0.1",
    username VARCHAR(255) NOT NULL DEFAULT "user",
    passwd VARCHAR(255) NOT NULL DEFAULT "changeme",
    topLevelDomain VARCHAR(255) NOT NULL DEFAULT "example",
    secondLevelDomain VARCHAR(255) NOT NULL DEFAULT "com",
    baseGroup VARCHAR(255) NOT NULL DEFAULT "dba",
    baseGroupOU VARCHAR(255) NOT NULL DEFAULT "groups",
    adminGroup VARCHAR(255) NOT NULL DEFAULT "sudodba",
    adminGroupOU VARCHAR(255) NOT NULL DEFAULT "groups",
    isDefault BOOLEAN NOT NULL DEFAULT FALSE,
    timeOutInSecs INT DEFAULT 3
);

-- create config here as we only want one row. Just to avoid any mistakes on the application
-- INSERT INTO config (id) VALUES (1);
-- This is for local testing only
-- INSERT INTO ldap_config(connectionStr, username, passwd, topLevelDomain)
--     VALUES("ldap://localhost:10389", "admin", "GoodNewsEveryone", "planetexpress");
INSERT INTO
    ldap_config(connectionStr, topLevelDomain)
VALUES
    ("ldap://localhost:10389", "planetexpress");

-- Adding tls to ldap connection
ALTER TABLE
    ldap_config
ADD
    COLUMN enableTLS BOOLEAN NOT NULL DEFAULT FALSE,
ADD
    COLUMN verifyCert BOOLEAN NOT NULL DEFAULT FALSE,
ADD
    COLUMN cert VARCHAR(255) NULL;