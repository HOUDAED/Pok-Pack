CREATE DATABASE IF NOT EXISTS pokepack;
USE pokepack;
CREATE
 IF NOT EXISTS TABLE `USER` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
    ,UNIQUE KEY uniq_username (username)
    ,UNIQUE KEY uniq_email (email)
);