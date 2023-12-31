-- Creation of the coursehub_api schema
CREATE SCHEMA coursehub_api;

-- Creation of the course table
CREATE TABLE coursehub_api.course
(
    id                SERIAL PRIMARY KEY,
    description       VARCHAR(50) NOT NULL,
    outline           TEXT        NOT NULL,
    registration_date TIMESTAMP DEFAULT timezone('BRT'::text, now())
);

-- Creation of a unique index on the 'id' column
CREATE UNIQUE INDEX coursehub_api_course_id_uindex ON coursehub_api.course (id);

-- Inserting sample data into the course table
INSERT INTO coursehub_api.course (description, outline, registration_date)
VALUES ('Matematica', 'Algebra', NOW()),
       ('Fisica', 'Ondas', NOW()),
       ('Química', 'Química Orgânica', NOW());

-- Creation of the student table
CREATE TABLE coursehub_api.student
(
    student_code SERIAL PRIMARY KEY,
    student_name VARCHAR(50)
);

-- Creation of the course_student table
CREATE TABLE coursehub_api.course_student
(
    registration_code SERIAL PRIMARY KEY,
    student_code      INT REFERENCES coursehub_api.student (student_code),
    course_code       INT REFERENCES coursehub_api.course (id)
);

-- Inserting sample data into the student and course_student tables
INSERT INTO coursehub_api.student (student_name)
VALUES ('Felipe Bravo'),
       ('Maria Silva'),
       ('João Santos');

INSERT INTO coursehub_api.course_student (student_code, course_code)
VALUES (1, 1),
       (2, 1),
       (2, 2),
       (3, 3);
