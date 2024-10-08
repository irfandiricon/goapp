# go-fiber (Example App GO)

## Introduction

go-fiber is a web application built with Go and Fiber that provides user authentication, profile management, and more. It uses MySQL for the database and JWT for authentication.

## Features

- User registration
- User login
- JWT-based authentication
- User profile management
- CRUD operations on users
- Pagination for listing users

## Project Structure

go-fiber/
├── controller/
│   ├── auth_controller.go
│   ├── user_controller.go
├── database/
│   ├── database.go
├── middleware/
│   ├── auth_middleware.go
├── model/
│   ├── user.go
├── routes/
│   ├── routes.go
├── .env
├── go.mod
├── go.sum
├── main.go
└── README.md


## Installation
- go mod download
- go run main.go


### Prerequisites

- Go 1.16+
- MySQL
- Git

### Clone the Repository

```sh
git clone https://github.com/irfandiricon/goapp.git
cd go-fiber
