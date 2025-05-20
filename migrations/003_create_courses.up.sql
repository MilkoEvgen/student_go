CREATE TABLE IF NOT EXISTS courses
(
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL UNIQUE,
    teacher_id BIGINT REFERENCES teachers (id) ON DELETE SET NULL
);