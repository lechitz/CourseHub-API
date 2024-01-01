CREATE SCHEMA coursehub_api

    CREATE TABLE course
    (
        id          SERIAL      NOT NULL
            CONSTRAINT coursehub_api_course_pkey PRIMARY KEY,
        description VARCHAR(50) NOT NULL,
        outline     TEXT        NOT NULL,
        registration_date  timestamp default timezone('BRT'::text, now())
    )

    CREATE UNIQUE INDEX coursehub_api_course_id_uindex ON course (id)

    CREATE TABLE student
    (
        id   SERIAL      NOT NULL
            CONSTRAINT coursehub_api_student_pkey PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        registration_date  timestamp default timezone('BRT'::text, now())
    )

    CREATE UNIQUE INDEX coursehub_api_student_id_uindex ON student (id)

    CREATE TABLE course_student
    (
        id            SERIAL NOT NULL
            CONSTRAINT coursehub_api_course_student_pkey PRIMARY KEY,
        fk_id_student INT    NOT NULL,
        fk_id_course  INT    NOT NULL
    )

    CREATE UNIQUE INDEX coursehub_api_course_student_id_uindex ON course_student (id);


    INSERT INTO coursehub_api.course (description, outline)
    VALUES ('Matematica', 'Algebra'),
           ('Fisica', 'Ondas'),
           ('Química', 'Química Orgânica');

    INSERT INTO coursehub_api.student (name)
    VALUES ('Felipe Bravo'),
           ('Maria Silva'),
           ('João Santos');

    INSERT INTO coursehub_api.course_student (fk_id_student, fk_id_course)
    VALUES (1, 1),
           (2, 1),
           (2, 2),
           (3, 3);