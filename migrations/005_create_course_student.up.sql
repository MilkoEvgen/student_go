CREATE TABLE IF NOT EXISTS course_student
(
    course_id  BIGINT REFERENCES courses (id) ON DELETE CASCADE,
    student_id BIGINT REFERENCES students (id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, student_id)
);