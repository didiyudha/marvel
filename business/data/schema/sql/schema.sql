-- version: 1.1
-- Description: Create table characters
create table "characters" (
    id int primary key,
    name varchar(255) not null default '',
    description text
);