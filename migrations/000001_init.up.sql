CREATE TABLE department(
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INTEGER NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_department_parent 
        FOREIGN KEY (parent_id) 
        REFERENCES department(id) ON DELETE CASCADE,
    CONSTRAINT department_name_parent_unique 
        UNIQUE (name, parent_id)
);

CREATE TABLE employee(
    id SERIAL PRIMARY KEY,
    department_id INTEGER NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_employee_department 
        FOREIGN KEY (department_id) 
        REFERENCES department(id) ON DELETE CASCADE
);