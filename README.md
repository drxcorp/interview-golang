# E-Commerce Backend API

A simple REST API backend for an e-commerce system built with Go 1.25.5.

## Features

- User management (create, read, update, delete, activate/deactivate)
- Product catalog management
- Order processing and tracking
- PostgreSQL database with Docker
- Database migrations using Goose

## Prerequisites

- Go 1.25.5 or later
- Docker and Docker Compose

## Quick Start

The easiest way to get started:

```bash
make start
```

This will:
1. Start PostgreSQL in Docker
2. Wait for it to be ready
3. Run the application (migrations run automatically on startup)

The server will start on `http://localhost:8080`

## Setup

1. Clone this repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   make start
   ```

## Available Make Commands

- `make start` - Start Docker and run the app (migrations run automatically)
- `make docker-up` - Start PostgreSQL in Docker
- `make docker-down` - Stop PostgreSQL Docker container
- `make run` - Run the application (migrations run automatically)

## Database Migrations

Migrations are embedded in the application and run automatically on startup. Migration files are located in the `database/migrations/` directory.

The application uses Goose as a Go library, so no separate installation is needed. All dependencies are managed through `go.mod`.

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint that returns server status, database connectivity, and user count

### Users
- `GET /users?search=<query>` - Get all users or search users
- `POST /users/create` - Create a new user
- `DELETE /users/delete?id=<id>` - Delete a user
- `GET /api/users?id=<id>` - Get user by ID
- `PUT /api/users` - Update user
- `GET /api/users/list?page=<page>&limit=<limit>` - List users with pagination
- `GET /api/users/search?q=<query>` - Search users
- `POST /api/users/activate?id=<id>` - Activate user
- `POST /api/users/deactivate?id=<id>` - Deactivate user

### Products
- `GET /products` - Get all products
- `POST /api/products` - Create a product
- `GET /api/products?id=<id>` - Get product by ID
- `PUT /api/products/stock?id=<id>&stock=<amount>` - Update product stock
- `DELETE /api/products?id=<id>` - Delete product
- `GET /api/products/all?min_price=<min>&max_price=<max>` - Get products with price filter
- `PUT /api/products/bulk-price?percentage=<percentage>` - Bulk update prices

### Orders
- `GET /orders?user_id=<id>` - Get all orders or orders by user
- `POST /orders/create` - Create a new order
- `POST /api/orders` - Create an order
- `GET /api/orders?id=<id>` - Get order by ID
- `GET /api/orders/user?user_id=<id>` - Get user's orders
- `PUT /api/orders/status?id=<id>&status=<status>` - Update order status
- `PUT /api/orders/cancel?id=<id>` - Cancel order
- `GET /api/orders/stats` - Get order statistics

## Environment Variables

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: testdb)
- `SERVER_PORT` - Server port (default: 8080)
- `APP_ENV` - Application environment (default: development)
- `LOG_LEVEL` - Log level (default: info)

## Default Data

The migration seeds the database with default data:
- User: John Doe (admin) - john@example.com / password123
- Product: Laptop - $999.99 (10 in stock)

## Docker Compose

PostgreSQL runs in a Docker container. Configuration is in `docker-compose.yml`. Data is persisted in a Docker volume.

To view logs:
```bash
docker-compose logs -f postgres
```

To connect to PostgreSQL:
```bash
docker exec -it interview_postgres psql -U postgres -d testdb
```
