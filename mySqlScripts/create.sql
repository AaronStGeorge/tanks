CREATE DATABASE IF NOT EXISTS Tanks;

USE Tanks;

CREATE TABLE Users (
id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
Username VARCHAR(60) NOT NULL UNIQUE,
Password VARCHAR(60) NOT NULL);


CREATE TABLE Friends (
user_id_1 INT NOT NULL,
user_id_2 INT NOT NULL);
