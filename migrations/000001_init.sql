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

INSERT INTO
    users (id, username, fio, PASSWORD)
VALUES (0, 'DELETED', 'DELETED', '');

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    author_id INT NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    title TEXT NOT NULL,
    preview_picture_url TEXT DEFAULT ''
);

INSERT INTO courses (id, author_id, title) VALUES (0, 0, 'DELETED');

CREATE TABLE course_lessons (
    id SERIAL PRIMARY KEY,
    lesson_number INT NOT NULL,
    course_id INT NOT NULL REFERENCES courses (id) ON UPDATE CASCADE ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    lesson_content TEXT,
    video_url TEXT DEFAULT ''
);

INSERT INTO
    course_lessons (
        id,
        lesson_number,
        course_id,
        title
    )
VALUES (0, 0, 0, 'DELETED');

CREATE TABLE lesson_attachments (
    id SERIAL PRIMARY KEY,
    attachment_type TEXT NOT NULL,
    lesson_id INT NOT NULL DEFAULT 0 REFERENCES course_lessons (id) ON UPDATE CASCADE ON DELETE SET DEFAULT,
    pretty_name TEXT,
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
    created_at DATE
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