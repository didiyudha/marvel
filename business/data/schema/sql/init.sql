CREATE DATABASE marvel;

create table "characters" (
    id int primary key,
    name varchar(255) not null default '',
    description text
);