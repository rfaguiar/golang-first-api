CREATE TABLE person(
     id SERIAL PRIMARY KEY,
     name VARCHAR (50) UNIQUE NOT NULL,
     age SMALLINT NOT NULL
);

INSERT INTO person(name, age)
VALUES('Jorge', 20);

INSERT INTO person(name, age)
VALUES('Jhon', 33);