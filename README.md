# Asset Management API

A RESTful API service for managing assets built with Go, Gin, and SQLite.

## Documentation

The API documentation is available at `URL_ADDRESS:YOUR_PORT/swagger/index.html`

## Features

- CRUD operations for assets
- Pagination support
- Sorting and ordering
- SQLite database
- Swagger documentation
- Unit tests with high coverage

## Prerequisites

- Go 1.24 or higher
- SQLite
- Docker (optional)

## Running Locally

1. Clone the repository

```bash
git clone git@github.com:mnasruls/assest-api-go.git
cd assets-api-go
```

2. Set up the environment variables

```bash
cp .env.example .env
```

3. Start the SQLite database

```bash
touch assets.db
```

4. Build and run the application

```bash
go run /cmd/main.go
```

5. The API will be available at `URL_ADDRESS:YOUR_PORT`

## Running with Docker

1. Build the Docker image

```bash
docker build -t assets-api .
```

2. Run the Docker container

```bash
docker run -p 9123:9123 assets-api
```

3. The API will be available at `URL_ADDRESS:9123`

## Running Tests

To run the unit tests, use the following command:

```bash
go test -v -coverprofile=coverage.out ./services
go tool cover -html=coverage.out
```
