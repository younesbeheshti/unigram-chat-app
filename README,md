# Chat Application Backend

A robust and scalable backend service for a real-time chat application built with Go.

## Features

- Real-time messaging using WebSocket
- User authentication and authorization
- PostgreSQL database integration
- RESTful API endpoints
- Secure password hashing
- JWT-based authentication
- Environment-based configuration

## Tech Stack

- **Language:** Go 1.23.5
- **Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **ORM:** GORM
- **WebSocket:** Gorilla WebSocket
- **Authentication:** JWT
- **Configuration:** godotenv

## Project Structure

```
.
├── config/         # Configuration files and database setup
├── handlers/       # HTTP request handlers
├── middleware/     # Custom middleware components
├── migrations/     # Database migrations
├── models/         # Data models and structures
├── routes/         # API route definitions
├── services/       # Business logic
├── storage/        # Database operations
├── utils/          # Utility functions
├── ws/             # WebSocket handlers
├── .env            # Environment variables
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
└── main.go         # Application entry point
```

## Prerequisites

- Go 1.23.5 or higher
- PostgreSQL database
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/younesbeheshti/chatapp-backend.git
   cd chatapp-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the root directory with the following variables:
   ```
   DB_HOST=your_db_host
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   DB_PORT=your_db_port
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `127.0.0.1:15000`.

## API Endpoints

The API documentation will be available at `/api/docs` when the server is running.

## WebSocket

The application supports real-time communication through WebSocket connections. Connect to the WebSocket endpoint at `ws://127.0.0.1:15000/ws`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Younes Beheshti

## Acknowledgments

- Gorilla WebSocket for WebSocket implementation
- GORM for database operations
- JWT for authentication
