-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS departments
(
    id serial primary key,
    name varchar(200) not null,
    parent_id bigint default null REFERENCES departments(id) on delete cascade,
    deleted_at timestamp,
    created_at timestamp not null default now(),
    CONSTRAINT department_name_not_empty CHECK (length(trim(name)) > 0)
);

CREATE UNIQUE INDEX parent_name_uq ON departments (COALESCE(parent_id, '-1'), name);
CREATE INDEX departments_parent_id_idx ON departments USING btree (parent_id);

-- Table: employee

-- DROP TABLE IF EXISTS employee;

CREATE TABLE IF NOT EXISTS employees
(
    id serial primary key,
    department_id bigint not null REFERENCES departments(id) on delete cascade,
    full_name varchar(200) not null,
    position varchar(200) not null,
    hired_at date,
    created_at timestamp not null default now(),
    CONSTRAINT employee_full_name_not_empty CHECK (length(trim(full_name)) > 0),
    CONSTRAINT employee_position_not_empty CHECK (length(trim(position)) > 0)
);

CREATE INDEX employees_department_id_idx ON employees USING btree (department_id);

ALTER TABLE IF EXISTS departments OWNER to postgres;
ALTER TABLE IF EXISTS employees OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
-- +goose StatementEnd
