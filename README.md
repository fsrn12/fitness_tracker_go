# Fitness Tracker Api in Go, PostgreSQL, Chi and Docker

# 🏋️ Fitness Tracker API

An authenticated, high-performance RESTful API built with **Go** and **PostgreSQL** following a **Clean Architecture** pattern. It provides robust CRUD operations for tracking user workouts.

---

## ✨ Key Features

- **Robust Authentication:** Secure user sign-up and sign-in using **JWTs** managed by dedicated token handlers and middleware.
- **Secure Middleware:** A custom authentication middleware that protects all private routes and validates incoming requests efficiently.
- **CRUD Functionality:** Full lifecycle management (Create, Read, Update, Delete) for user-specific workout data.
- **Clean Architecture:** Clear separation of concerns into API handlers, business logic, and a data access layer (`store`).
- **Database:** Persistent storage using **PostgreSQL** managed via **Docker Compose**.
- **Performant:** Built with standard Go libraries, the Chi router, and the high-efficiency `pgx` driver for raw SQL queries.

## 🛠️ Tech Stack

| Component            | Technology     | Role                                                  |
| :------------------- | :------------- | :---------------------------------------------------- |
| **Language**         | Go (Golang)    | Primary backend language.                             |
| **Routing**          | Chi            | High-performance, minimal, and idiomatic Go router.   |
| **Database**         | PostgreSQL     | Robust and reliable relational database.              |
| **Driver**           | `pgx`          | Modern, high-performance PostgreSQL driver.           |
| **Configuration**    | `dotenv`       | Secure handling of environment variables and secrets. |
| **Containerization** | Docker Compose | Manages the local PostgreSQL database environment.    |

---

## 📂 Architecture & Directory Structure

This project adopts a **Clean Architecture** approach to maximize scalability, testability, and maintainability.
