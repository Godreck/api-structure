-- +goose Up
-- +goose StatementBegin
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id BIGINT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_departments_parent
        FOREIGN KEY (parent_id)
        REFERENCES departments(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE UNIQUE INDEX ux_departments_parent_name
    ON departments (parent_id, name);

CREATE INDEX ix_departments_parent_id
    ON departments (parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd
