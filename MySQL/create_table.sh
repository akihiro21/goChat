#!/bin/sh

CMD_MYSQL="mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}"
$CMD_MYSQL -e "create table user (
    id int(10)  AUTO_INCREMENT NOT NULL primary key,
    name varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    room int(10)
    );"
$CMD_MYSQL -e "create table room (
    id int(10) AUTO_INCREMENT NOT NULL primary key,
    name varchar(50) NOT NULL,
    date datetime,
    userId1 int(10),
    userId2 int(10)
    );"

$CMD_MYSQL -e "create table message (
    id int(10) AUTO_INCREMENT NOT NULL primary key,
    message varchar(255) NOT NULL,
    room varchar(50) NOT NULL,
    user varchar(50) NOT NULL
    );"