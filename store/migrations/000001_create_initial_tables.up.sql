-- -----------------------------------------------------
-- Schema
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS "gapp";
SET search_path TO "$user","gapp","public";

-- -----------------------------------------------------
-- enable UUID extension
-- -----------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -----------------------------------------------------
-- Create tables
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS gapp.rooms
(
    id           uuid DEFAULT  uuid_generate_v4() NOT NULL ,
    availability TEXT,
    created      TIMESTAMP,
    modified     TIMESTAMP
)