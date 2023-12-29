create schema coursehub_api

-- auto-generated definition
    create table course
    (
        id serial not null
            constraint coursehub_api_course_pkey primary key,
        description varchar(50) not null,
        outline varchar(255) not null,
        registration_date timestamp default timezone('BRT'::text, now())
    )

    create
        unique index coursehub_api_course_id_uindex
        on cliente (id);

-- Create default accounts inserts
INSERT INTO coursehub_api.course (description, outline, registration_date)
VALUES ('matematica', 'equações de 1 grau', now());