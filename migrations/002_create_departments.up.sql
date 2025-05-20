CREATE TABLE IF NOT EXISTS departments
(
    id                    BIGSERIAL PRIMARY KEY,
    name                  TEXT NOT NULL,
    head_of_department_id BIGINT references teachers (id) ON DELETE SET NULL
);