Project - AMPL TECHNOLOGY Task Management Test assignment


Overview

A RESTful API built with the Gin framework for managing tasks. It supports creating, retrieving, updating, and deleting tasks, with JWT-based authentication.
Setup

Requirements

    Go 1.20+
    MySql

Steps

    Clone the repository:

git clone https://github.com/VandanaSingh2209/task-management-api


cd task-management-api

Install dependencies:

go mod tidy



project-root/
│
├── Controllers/
│   └── (controller files like GetTaskList, CreateTask, etc.)
├── Database/
│   └── (database connection file like Database.go)
├── Middleware/
│   └── (middleware files like Auth.go)
├── Models/
│   └── (model definitions like AmplTaskList.go)
├── Routes/
│   └── (routes setup like routes.go)
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── .env

Configure .env with database credentials:

DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=AmplTaskList

Run the application:

    go run .

-- Create the database
CREATE DATABASE AmplTest;

-- Use the database
USE AmplTest;

-- Create the table
CREATE TABLE AmplTaskList (
    ID INT AUTO_INCREMENT PRIMARY KEY, -- Auto-generated unique ID
    Title VARCHAR(255) NOT NULL, -- Title (string, required)
    Description TEXT NOT NULL, -- Description (string, required)
    Status ENUM('pending', 'in-progress', 'completed') NOT NULL DEFAULT 'pending', -- Status (enum)
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Auto-generated timestamp for creation
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Auto-updated timestamp
);


Middleware

1. JWT Authentication

All routes under /protected require a valid JWT token. The middleware verifies the Authorization header.

Tokens are generated using the JWTTokenGenerate function and signed with a secret key.

Token structure:

    sub: User ID
    exp: Expiration time (1 hour)
    iss: Issuer

Error Handling:

    If a token is missing, invalid, or expired, the middleware returns a 401 Unauthorized error.



2. Rate Limiting

Rate limiting is applied to all /protected routes. By default, the middleware limits requests to:
5 requests per 10 seconds

    The RateLimitMiddleware restricts the number of requests a client can make to protected routes.
    Configuration:
        Requests per second (RPS): Adjustable.
        Burst size: Number of requests allowed at once before rate-limiting starts.
    Default Behavior:
        Too Many Requests error is returned if the rate limit is exceeded.



API Endpoints
1. Get Task List

    GET hostname/public/taskslist?page=1&limit=10
    Response: Paginated list of tasks.

2. Get Task

    GET hostname/private/task?id={id}
    Response: Task details.

3. Create Task

    POST hostname/private/create-task
    Body:

    {"title": "New Task", "description": "Task Description", "status": "pending"}

    Response: Task created.

4. Update Task

    PUT hostname/private/tasks/{id}
    Body:

    {"title": "Updated Task", "description": "Updated Description", "status": "in-progress"}

    Response: Task updated.

5. Delete Task

    DELETE hostname/private/tasks/{id}
    Response: Task deleted.

6. Login

    POST hostname/public/login
    Body:

    {"userid": 123}

    Response: JWT token.

Testing

    Install test dependencies:

go get github.com/stretchr/testify github.com/DATA-DOG/go-sqlmock

Run tests:

    go test ./...

