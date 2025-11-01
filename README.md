# Microservices Project

[![CircleCI](https://circleci.com/gh/igomez10/microservices/tree/mainline.svg?style=svg)](https://circleci.com/gh/igomez10/microservices/tree/mainline)

A collection of microservices built with Go, demonstrating modern cloud-native development practices including API design, observability, and deployment automation.

## 📋 Table of Contents

- [Overview](#overview)
- [Services](#services)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Development](#development)
- [Deployment](#deployment)
- [Observability](#observability)
- [Contributing](#contributing)
- [License](#license)

## 🔍 Overview

This repository contains multiple microservices that demonstrate best practices in modern microservice architecture, including:

- RESTful API design using OpenAPI 3.1 specifications
- OAuth2 authentication and authorization
- Comprehensive observability stack (metrics, logs, traces)
- Contract testing with Pact
- Infrastructure as Code
- CI/CD pipelines

## 🚀 Services

### Socialapp

A full-featured social networking API with user management, comments, followers, and role-based access control.

- **Location**: `socialapp/`
- **Base URL**: https://socialapp.gomezignacio.com
- **API Documentation**: [socialapp/README.md](socialapp/README.md)
- **Features**:
  - User management with authentication
  - Comments and user feeds
  - Follow/unfollow functionality
  - Role-based access control (RBAC)
  - OAuth2 security with scopes
  - URL shortening integration

### URL Shortener

A high-performance URL shortening service with analytics and OAuth2 protection.

- **Location**: `urlshortener/`
- **Base URL**: https://urlshortener.gomezignacio.com
- **API Documentation**: [urlshortener/README.md](urlshortener/README.md)
- **Features**:
  - Create and manage short URLs
  - URL metadata tracking
  - OAuth2 protected endpoints
  - Analytics and usage tracking

### Bitcoin Price Service

A service for tracking and displaying Bitcoin price information using the Buffalo framework.

- **Location**: `bitcoinprice/`
- **Framework**: Buffalo (Go)
- **Features**:
  - Bitcoin price tracking
  - Real-time updates
  - Historical data

## 🛠 Technology Stack

### Core Technologies

- **Language**: Go 1.23
- **API Framework**: Chi router
- **API Specification**: OpenAPI 3.1.0
- **Database**: PostgreSQL (with pgx driver)
- **Caching**: Redis
- **Authentication**: OAuth2 with JWT

### Code Generation

- **OpenAPI Generator**: Automated client/server code generation
- **SQLC**: Type-safe SQL code generation

### Observability Stack

- **Metrics**: Prometheus, Grafana Mimir
- **Tracing**: OpenTelemetry, Tempo, Jaeger
- **Logging**: Logstash
- **Profiling**: Pyroscope (continuous profiling)
- **APM**: New Relic
- **Visualization**: Grafana

### Infrastructure & Deployment

- **Containerization**: Docker, Docker Compose
- **Reverse Proxy**: Traefik, Caddy
- **CI/CD**: CircleCI, GitHub Actions
- **Infrastructure as Code**: Terraform
- **Load Testing**: JMeter

### Testing

- **Unit Testing**: Go testing framework
- **Contract Testing**: Pact
- **Integration Testing**: Custom test suites
- **API Testing**: Postman collections

### Development Tools

- **Change Data Capture**: Debezium
- **Hot Reload**: Reflex
- **Code Quality**: Go fmt, goimports

## 🏗 Architecture

The project follows a microservices architecture with the following patterns:

- **API-First Design**: All services are defined using OpenAPI specifications
- **OAuth2 Security**: Centralized authentication and authorization
- **Database per Service**: Each microservice has its own PostgreSQL database
- **Event-Driven Communication**: Change Data Capture with Debezium
- **Comprehensive Observability**: Full metrics, logs, and traces for all services

### Service Communication

- RESTful APIs for synchronous communication
- OAuth2 for secure inter-service communication
- Redis for caching and session management

## 🚦 Getting Started

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- PostgreSQL 9.6+
- Redis
- Make

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/igomez10/microservices.git
   cd microservices
   ```

2. **Start a service with Docker Compose**
   
   For Socialapp:
   ```bash
   cd socialapp
   docker compose up -d
   ```
   
   For URL Shortener:
   ```bash
   cd urlshortener
   docker compose up -d
   ```

3. **Access the services**
   - Socialapp: http://localhost (routed via Traefik on port 80)
   - URL Shortener: http://localhost:8089

### Local Development Setup

Each service has its own development setup. See individual service READMEs for detailed instructions.

**Socialapp:**
```bash
cd socialapp
make generate-openapi  # Generate API code
make sqlc-generate     # Generate SQL code
make build             # Build the service
make test              # Run tests
```

**URL Shortener:**
```bash
cd urlshortener
make generate-openapi  # Generate API code
make sqlc-generate     # Generate SQL code
make build             # Build the service
make test              # Run tests
```

**Bitcoin Price:**
```bash
cd bitcoinprice
buffalo dev            # Start development server
buffalo test           # Run tests
```

## 💻 Development

### API Code Generation

Both socialapp and urlshortener use OpenAPI Generator to create:
- Server stubs
- Client SDKs (Go, Rust, etc.)
- API documentation (Markdown)
- CLI tools (Bash)
- Postman collections
- JMeter test plans

### Database Migrations

Services use SQLC for type-safe SQL code generation:
```bash
make sqlc-generate
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .
```

## 🚀 Deployment

### Automated Deployment

The project uses GitHub Actions for continuous deployment:

- **Socialapp**: Deploys on push to `mainline` branch when files in `socialapp/` change
- **URL Shortener**: Deploys on push to `mainline` branch when files in `urlshortener/` change

### Manual Deployment

Each service includes a deployment script:
```bash
cd <service-directory>
make deploy
```

### Infrastructure

Infrastructure configuration is managed with Terraform:
```bash
cd socialapp/infrastructure
terraform init
terraform plan
terraform apply
```

## 📊 Observability

The project includes a comprehensive observability stack:

### Metrics
- Prometheus for metrics collection
- Grafana Mimir for long-term metrics storage
- Grafana for visualization
- Custom dashboards for each service

### Tracing
- OpenTelemetry for trace instrumentation
- Tempo for trace storage
- Jaeger for trace visualization
- Distributed tracing across services

### Logging
- Structured logging with zerolog
- Logstash for log aggregation

### Profiling
- Pyroscope for continuous profiling
- CPU and memory profiling
- Performance analysis

### Monitoring Setup

Start the observability stack:
```bash
cd socialapp
docker compose up -d prometheus grafana tempo pyroscope
```

Access dashboards:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Ensure all tests pass before submitting PR

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Ignacio Gomez**
- Email: ignacio.gomez.arboleda@gmail.com
- GitHub: [@igomez10](https://github.com/igomez10)

## 🔗 Links

- [Socialapp API](https://socialapp.gomezignacio.com)
- [URL Shortener API](https://urlshortener.gomezignacio.com)
- [CircleCI Builds](https://circleci.com/gh/igomez10/microservices)

---

**Note**: This is an active development project. Some features may be work in progress.
