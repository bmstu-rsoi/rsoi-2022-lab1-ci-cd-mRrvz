-- file: 10-create-user-and-db.sql
CREATE DATABASE persons;
CREATE ROLE program WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE persons TO program;
ALTER ROLE program WITH LOGIN;

\c persons;

CREATE TABLE persons (
    id serial PRIMARY KEY,
    name VARCHAR (64),
    age INT,
    address VARCHAR(128),
    work VARCHAR(128)
);

INSERT INTO persons(name, age, address, work) VALUES
    ('Dmitry Yakuba', 22, 'Yasenevo', 'IVI'),
    ('Pavel Perestoronin', 22, 'Yasenevo', 'Yandex'),
    ('Mikhail Nitenko', 22, 'Mauritius', 'Parallels'),
    ('Sergey Kononenko', 22, 'Astana', 'SberMarket');
