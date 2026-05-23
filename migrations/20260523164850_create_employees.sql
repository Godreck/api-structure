-- +goose Up
-- +goose StatementBegin
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_employees_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX ix_employees_department_id
    ON employees (department_id);

CREATE INDEX ix_employees_full_name
    ON employees (full_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
-- +goose StatementEnd
