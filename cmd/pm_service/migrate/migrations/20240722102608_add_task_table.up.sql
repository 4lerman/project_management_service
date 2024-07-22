CREATE TYPE task_type AS ENUM ('low', 'medium', 'high');

CREATE TYPE task_priority AS ENUM ('new', 'in_process', 'done');

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) NOT NULL,
    descript VARCHAR(255),
    taskType task_type NOT NULL,
    taskPriority task_priority NOT NULL,
    userId INT NOT NULL,
    projectId INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (projectId) REFERENCES projects(id)
);