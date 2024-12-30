# Nova

## Lightweight API Framework Based on Golang

![](https://img.shields.io/npm/l/vue.svg)
![dbtest](https://github.com/boyxp/nova/actions/workflows/go.yml/badge.svg)

[中文文档](README.CN.md)

This document describes a lightweight API framework based on Golang, designed to simplify the API development process and provide high performance. The following are the main features of this framework:

## Features

*   One-click project initialization, rapid development: Provides scaffolding tools that can quickly build project skeletons, reducing initial configuration work and allowing developers to focus on writing business logic.
*   Support for graceful restart: Restarts the server without interrupting existing requests, ensuring service continuity, which is very important for online services.
*   Automatic route registration: Automatically scans and registers routes without manual configuration, simplifying route management and improving development efficiency.
*   Coding-friendly ORM queries, no field mapping required, no need to define return structs: Provides easy-to-use ORM (Object-Relational Mapping) tools to simplify database operations, eliminating the need for manual field mapping. It can also flexibly handle query results without pre-defining return structs, improving development efficiency and flexibility.
*   Automatic request parameter validation and data type conversion: Automatically validates the validity of request parameters and performs data type conversion, reducing the amount of manual validation and conversion code, and improving code robustness.
*   No need to pass Context in requests: The framework handles Context internally, and developers do not need to explicitly pass Context in business logic, simplifying code writing.
*   Exception handling familiar to PHP developers: Provides an exception handling mechanism similar to PHP, making it easy for PHP developers to get started quickly and improve code maintainability.
*   Low-intrusion design, existing structs can be exposed as services with one line of code: Has extremely low intrusion into existing code. Existing structs can be exposed as API interfaces with just one line of code, facilitating rapid integration and transformation of existing projects.
*   Middleware support: Supports the use of middleware to handle pre- and post-request logic, such as logging, authentication, authorization, etc., improving code reusability and maintainability.
*   Loosely coupled modular design: Adopts a modular design, with independent modules that are easy to expand and maintain.

## Quick Start

1. **Initialize a New Go Module:**
   Create a new Go module for your project.

   ```bash
   go mod init api
   ```

2. **Create a Basic Application File:**
   Create a file named `hello.go` with the following content:

   ```go
   package main

   import (
       "github.com/boyxp/nova"
       "github.com/boyxp/nova/router"
   )

   // Define a struct for your controller
   type Hello struct {}

   // Define a method on your struct to handle requests
   func (h *Hello) Hi(name string) map[string]string {
       return map[string]string{"name": "hello " + name}
   }

   func main() {
       // Register the controller with the router
       router.Register(Hello{})
       // Start the Nova server on port 9800
       nova.Listen("9800").Run()
   }
   ```

3. **Install Dependencies:**
   Use `go mod tidy` to install the necessary dependencies.

   ```bash
   go mod tidy
   ```

4. **Run the Application:**
   Start your application using:

   ```bash
   go run hello.go
   ```

5. **Test the API:**
   Use `curl` to test the API endpoint:

   ```bash
   curl -X POST -d 'name=eve' 127.0.0.1:9800/hello/hi
   ```

## Project Mode

To initialize a complete project structure and manage the application, follow these additional steps:

### Initialize the Project:
   Run the initialization script to set up the project structure.

   ```bash
   curl -s https://raw.githubusercontent.com/boyxp/nova/master/init.sh | sh
   ```

### Process Management:
   Use the provided `manage.sh` script to manage the application process.

   - **Start Process:**
     ```bash
     sh manage.sh start
     ```

   - **Check Process Status:**
     ```bash
     sh manage.sh status
     ```

   - **Smooth Restart (build and restart without interrupting current requests):**
     ```bash
     sh manage.sh restart
     ```

   - **Reload Configuration (without build):**
     ```bash
     sh manage.sh reload
     ```

   - **Stop Process (after completing current requests):**
     ```bash
     sh manage.sh stop
     ```

   - **Upgrade Code and Restart:**
     ```bash
     sh manage.sh upgrade
     ```

### Creating a Controller

1. **Define a Controller:**
   Create a new controller in the `controller` directory and register it with the router.

   ```go
   package controller

   import "github.com/boyxp/nova/router"

   // Initialize and register the controller
   func init() {
       router.Register(Hello{})
   }

   // Define the controller struct
   type Hello struct {}

   // Define a method to handle requests
   func (h *Hello) Hi(name string) map[string]string {
       return map[string]string{"name": "hello " + name}
   }
   ```

2. **Restart the Application:**
   Restart the process to apply changes.

   ```bash
   sh manage.sh restart
   ```

3. **Test the API:**
   Use `curl` to test the new controller's endpoint:

   ```bash
   curl -X POST -d 'name=eve' 127.0.0.1:9800/hello/hi
   ```

### Additional Configuration

- **Cookie Settings:**
  Configure cookie settings in your environment variables (e.g., `cookie.HttpOnly`, `cookie.Secure`, `cookie.Path`, `cookie.Domain`, `cookie.MaxAge`).

- **Database Configuration:**
  Register database parameters and create models for database operations.

  ```go
  // Register database parameters
  database.Register("database", "test", "root:123456@tcp(localhost:3306)/test")

  // Create a model
  Goods := database.Model{"database.goods"}

  // Perform a query
  Goods.Where("1").Find()
  ```

