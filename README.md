# IU-K8s Backend API

A REST API server built with Go, using oapi-codegen and chi router, following OpenAPI 3.0 specification.

## Features

- 🚀 Fast HTTP server with chi router
- 📝 OpenAPI 3.0 specification-driven development
- 🔄 Auto-generated code with oapi-codegen
- 🏗️ Clean architecture with layered structure
- 🔧 Configuration via environment variables
- 📊 Health check endpoint
- 👥 User management CRUD operations

## Project Structure

```
.
├── cmd/
│   └── server/                 # Application entry point
│       └── main.go
├── internal/
│   ├── api/                    # Generated API code
│   │   └── generated.go
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── handlers/              # HTTP handlers
│   │   └── user_handler.go
│   ├── middleware/            # HTTP middleware
│   │   └── middleware.go
│   ├── repository/            # Data access layer
│   │   └── user_repository.go
│   └── service/               # Business logic layer
│       └── user_service.go
├── openapi.yaml               # OpenAPI specification
├── openapi_config.yaml        # oapi-codegen configuration
├── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24.5 or later
- Git

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd iu-k8s
```

2. Install dependencies:

```bash
go mod tidy
```

3. Generate API code (if needed):

```bash
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config openapi_config.yaml openapi.yaml
```

### Running the Server

#### Using default configuration:

```bash
go run cmd/server/main.go
```

#### Using environment variables:

```bash
export PORT=8080
export READ_TIMEOUT=10
export WRITE_TIMEOUT=10
export IDLE_TIMEOUT=120
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Environment Variables

| Variable        | Description              | Default |
| --------------- | ------------------------ | ------- |
| `PORT`          | Server port              | `8080`  |
| `READ_TIMEOUT`  | Read timeout in seconds  | `10`    |
| `WRITE_TIMEOUT` | Write timeout in seconds | `10`    |
| `IDLE_TIMEOUT`  | Idle timeout in seconds  | `120`   |

## API Endpoints

### Health Check

- `GET /health` - Returns server health status

### Users

- `GET /api/v1/users` - List users with pagination
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{userId}` - Get user by ID
- `PUT /api/v1/users/{userId}` - Update user
- `DELETE /api/v1/users/{userId}` - Delete user

### Example Usage

#### Health Check

```bash
curl http://localhost:8080/health
```

#### Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "John Doe"
  }'
```

#### List Users

```bash
curl "http://localhost:8080/api/v1/users?limit=10&offset=0"
```

#### Get User

```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

## Development

### Code Generation

The API code is generated from the OpenAPI specification. To regenerate:

```bash
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config openapi_config.yaml openapi.yaml
```

### Project Architecture

This project follows a clean architecture pattern:

1. **API Layer** (`internal/api/`) - Generated from OpenAPI spec
2. **Handler Layer** (`internal/handlers/`) - HTTP request/response handling
3. **Service Layer** (`internal/service/`) - Business logic
4. **Repository Layer** (`internal/repository/`) - Data access

### Adding New Endpoints

1. Update `openapi.yaml` with new endpoint specifications
2. Regenerate API code using oapi-codegen
3. Implement handlers in `internal/handlers/`
4. Add business logic in `internal/service/`
5. Update repository if needed

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Building

```bash
# Build for current platform
go build -o bin/server cmd/server/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

## Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

Build and run:

```bash
docker build -t iu-k8s-api .
docker run -p 8080:8080 iu-k8s-api
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
