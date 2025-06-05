# Go Event Management API

A RESTful API for managing events and user registrations, built with Go, Gin, and SQLite.

## Features

- User authentication with JWT
- Event creation and management
- Event registration system
- Swagger documentation
- Clean architecture
- SQLite database

## Prerequisites

- Go 1.24 or higher
- SQLite3

## Installation

1. Clone the repository:
```bash
git clone https://github.com/gojo-op/go-event-app.git
cd go-event-app
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (optional):
```bash
export SERVER_MODE=debug
export SERVER_PORT=8000
export LOG_LEVEL=info
export JWT_SECRET=your-256-bit-secret
```

4. Run the application:
```bash
go run main.go
```

Or with live reload using Air:
```bash
air
```

## Project Structure

```
go-event-app/
├── docs/                  # Swagger documentation
├── internal/             # Internal packages
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── model/          # Data models
│   ├── repository/     # Database operations
│   ├── service/        # Business logic
│   └── utils/          # Utility functions
├── events.db            # SQLite database
├── main.go             # Application entry point
├── go.mod              # Go module file
└── README.md           # This file
```

## Architecture

The application follows a clean, layered architecture:

1. **Handler Layer** (Presentation)
   - Handles HTTP requests and responses
   - Input validation
   - Response formatting
   - Uses the standardized response utilities

2. **Service Layer** (Business Logic)
   - Implements business rules
   - Coordinates between handlers and repositories
   - Handles data transformation
   - Manages transactions

3. **Repository Layer** (Data Access)
   - Database operations
   - Data persistence
   - Query execution
   - Raw data retrieval

4. **Model Layer** (Data Structures)
   - Defines data structures
   - Validation rules
   - Data transfer objects (DTOs)

## Control Flow

Here's how a typical request flows through the application:

1. **Request Entry**
   ```
   HTTP Request → Gin Router → Middleware → Handler
   ```
   - Request hits a route
   - Middleware processes request (auth, logging, etc.)
   - Handler receives the processed request

2. **Handler Processing**
   ```
   Handler → Input Validation → Service Call → Response Formatting
   ```
   - Validates input data
   - Calls appropriate service method
   - Formats response using utils.Response

3. **Service Processing**
   ```
   Service → Business Logic → Repository Calls → Data Transformation
   ```
   - Applies business rules
   - Makes necessary repository calls
   - Transforms data as needed

4. **Repository Processing**
   ```
   Repository → SQL Query → Database → Data Mapping
   ```
   - Executes SQL queries
   - Maps database results to models
   - Handles database errors

5. **Response Flow**
   ```
   Database → Repository → Service → Handler → HTTP Response
   ```
   - Data flows back up through layers
   - Each layer adds its processing
   - Final response is formatted and sent

## Standardized Response Format

All API responses follow a standard format:

```json
// Success Response
{
    "success": true,
    "message": "Operation successful",
    "data": {
        // Response data here
    }
}

// Error Response
{
    "success": false,
    "message": "Error description",
    "error": "Detailed error message"
}

// Paginated Response
{
    "success": true,
    "message": "Operation successful",
    "data": [...],
    "pagination": {
        "current_page": 1,
        "total_pages": 10,
        "total_records": 100,
        "page_size": 10
    }
}
```

## Database Schema

```sql
-- Users Table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Events Table
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    date DATETIME NOT NULL,
    capacity INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- Registrations Table
CREATE TABLE registrations (
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (user_id, event_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);
```

## API Endpoints

### User Endpoints
- `POST /users/register` - Register a new user
- `POST /users/login` - Login user
- `GET /users/me` - Get current user profile
- `PUT /users/me` - Update current user profile

### Event Endpoints
- `POST /events` - Create a new event
- `GET /events` - List all events
- `GET /events/:id` - Get event details
- `PUT /events/:id` - Update event
- `DELETE /events/:id` - Delete event
- `POST /events/:id/register` - Register for event
- `POST /events/:id/unregister` - Unregister from event
- `GET /events/:id/registrations` - Get registered users

## Authentication

The application uses JWT (JSON Web Tokens) for authentication:
1. User logs in with credentials
2. Server validates and returns a JWT
3. Client includes JWT in Authorization header
4. Protected routes verify JWT in middleware

## Setup and Running

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Access the API documentation:
   ```
   http://localhost:8080/swagger/index.html
   ```

## Error Handling

The application uses a centralized error handling approach:
1. Repository layer returns specific database errors
2. Service layer translates technical errors to business errors
3. Handler layer converts errors to appropriate HTTP responses
4. Response utilities ensure consistent error formatting

## Development Workflow

1. Define models in `internal/model`
2. Implement repository methods in `internal/repository`
3. Add business logic in `internal/service`
4. Create handlers in `internal/handler`
5. Add routes in `main.go`
6. Update Swagger documentation
7. Test endpoints

## Best Practices

1. **Separation of Concerns**
   - Each layer has a specific responsibility
   - Clear boundaries between layers
   - Dependency injection for better testing

2. **Error Handling**
   - Consistent error types
   - Proper error propagation
   - User-friendly error messages

3. **Security**
   - Password hashing
   - JWT authentication
   - Input validation
   - SQL injection prevention

4. **Code Organization**
   - Clear package structure
   - Consistent naming conventions
   - Proper documentation
   - Clean interfaces

5. **Database**
   - Prepared statements
   - Transaction support
   - Connection pooling
   - Proper indexing

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
