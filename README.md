# Questify
## Overview
Questify is a survey management system with robust role-based access control (RBAC). It provides functionality to create surveys, assign roles to users, and manage permissions for various operations. This document outlines the features, installation steps, and usage of the application.

---
# Features
## 1. User Authentication and Authorization
- JWT-based authentication.
- Role-based access control with permissions for fine-grained access.
- 
## 2. Survey Management
- Create, update, and delete surveys.
- Add and manage questions within surveys.
- Assign roles and permissions to users for specific surveys.

## 3. Roles and Permissions
- Define roles such as admin, editor, viewer, etc.
- Assign permissions to roles to restrict access to specific actions.

## 4. API Endpoints
-User management: Registration, login, and role assignments.
-Survey management: Create, edit, delete, and view surveys.
-Role management: Assign and revoke roles with associated permissions.

# How To Run
1. copy the code from click[here](https://github.com/hesamhme/Questify/blob/development/config-sample.txt) sample file and past to config.yaml
2. run the docker compose
```
docker-compose up -d
```
3. run the project
```
go mod tidy
go run cmd/api/main.go
```

## Tech Stack:

- **Go** (v1.22+): Primary programming language of the project. <span><img src="https://img.shields.io/badge/Golang-1.23-blue" /></span>
- **PostgreSQL**: Database used for storing data. <span><img src="https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white" /></span>
- **GORM**: ORM used for database connection.
- **Fiber**: The web framework used for REST API.
- **Docker**: The contrinize app.<span><img src="https://img.shields.io/badge/Docker-2CA5E0?style=flat&logo=docker&logoColor=white" /></span>

## ERD (Entity Relationship Diagram)

- click [here](https://dbdiagram.io/d/quera_p2-6744b17ce9daa85acaa836e1) to check the ERD.
- click [here](https://restless-star-600413.postman.co/workspace/Quera~3fa01880-1528-4b61-8a86-e910c17ad706/collection/24632383-30a3ba92-d18e-4d5d-aac6-1e9d0c1829e5?action=share&source=copy-link&creator=24632383) to Postman collections. 

![quera_p2](https://github.com/user-attachments/assets/991b5261-1193-412c-9812-2d73c34e9821)
