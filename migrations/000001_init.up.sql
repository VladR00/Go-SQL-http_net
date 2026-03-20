CREATE TABLE Department(
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INTEGER NULL,
    created_at TIMESTAMP,
    CONSTRAINT fk_department_parent 
        FOREIGN KEY (parent_id) 
        REFERENCES Department(id),
    CONSTRAINT department_name_parent_unique 
    UNIQUE (name, parent_id)
);

CREATE TABLE Employee(
    id SERIAL PRIMARY KEY,
    department_id INTEGER,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMP,
    CONSTRAINT fk_employee_department 
        FOREIGN KEY (department_id) 
        REFERENCES Department(id)
);