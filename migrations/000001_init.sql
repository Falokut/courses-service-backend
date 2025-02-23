-- +goose Up
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO roles (id, name) VALUES (0, 'none');

INSERT INTO roles (name) VALUES ('admin'), ('student'), ('lector');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    fio TEXT NOT NULL,
    PASSWORD TEXT NOT NULL,
    role_id INT NOT NULL DEFAULT 0 REFERENCES roles (id) ON UPDATE CASCADE ON DELETE SET DEFAULT
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name TEXT NOT NULL,
    preview_picture_url TEXT DEFAULT ''
);

CREATE TABLE course_lessons (
    id SERIAL PRIMARY KEY,
    course_id INT NOT NULL REFERENCES courses (id) ON UPDATE CASCADE ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    description TEXT
);

CREATE TABLE lesson_attachments (
    id SERIAL PRIMARY KEY,
    attachment_type TEXT NOT NULL,
    lesson_id INT NOT NULL REFERENCES course_lessons (id) ON UPDATE CASCADE ON DELETE CASCADE,
    url TEXT NOT NULL
);

CREATE TABLE courses_registration (
    course_id INT NOT NULL REFERENCES courses (id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (course_id, user_id)
);

CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    course_id INT NOT NULL REFERENCES courses (id) ON UPDATE CASCADE ON DELETE CASCADE,
    lector_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    assignment_text TEXT,
    assignment_at DATE
);

CREATE TABLE submitted_assignments (
    student_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    assignment_id INT NOT NULL REFERENCES assignments (id) ON UPDATE CASCADE ON DELETE CASCADE,
    grade INT NOT NULL CHECK (
        grade > 0
        AND grade < 6
    ),
    submitted_at DATE,
    attachments_info JSON DEFAULT '{}',
    PRIMARY KEY (student_id, assignment_id)
);

-- +goose Down
DROP TABLE submitted_assignments;

DROP TABLE assignments;

DROP TABLE lesson_attachments;

DROP TABLE course_lessons;

DROP TABLE courses;

DROP TABLE sessions;

DROP TABLE users;

DROP TABLE roles;