create database arenda;

create table city
(
    id   SERIAL primary key,
    name text not null
);

create table district
(
    id   serial primary key,
    name text not null
);

create table street
(
    id   serial primary key,
    name text not null
);

create table house
(
    id          serial primary key,
    city_id     int references city (id) on update cascade,
    district_id int references district (id) on update cascade,
    street_id   int references street (id) on update cascade,
    "number"    int not null,
    literal     text,
    floor_count int not null
);

create table flat
(
    id         serial primary key,
    house_id   int references house (id) on update cascade,
    "number"   int            not null,
    "cost"     numeric(20, 2) not null,
    "space"    numeric(20, 2) not null,
    "floor"    smallint       not null,
    room_count smallint       not null
);