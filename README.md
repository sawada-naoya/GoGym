# GoGym Next - Gym Search & Review Platform

A modern gym search and review platform built with Go backend and Next.js frontend, featuring location-based search, user reviews, and photo uploads.

## 🏗️ Architecture

This is a monorepo containing:

- **Backend**: Go 1.22+ with Echo framework, Clean Architecture
- **Frontend**: Next.js 14+ with App Router, React Server Components
- **Database**: MySQL 8.0 with spatial support
- **Storage**: S3-compatible (MinIO for dev)
- **Cache**: Redis
- **Containerization**: Docker Compose

## 🚀 Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Task](https://taskfile.dev/#/installation) - task runner
- [Go 1.22+](https://golang.org/dl/) (for local development)
- [Node.js 18+](https://nodejs.org/) and [pnpm](https://pnpm.io/) (for local development)

### Setup

1. **Clone and setup environment**
   ```bash
   git clone <repository-url>
   cd GoGym
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start development environment**
   ```bash
   task dev
   ```

   This command will:
   - Generate code from OpenAPI spec
   - Start all services (MySQL, Redis, MinIO, API, Web)
   - Run database migrations
   - Seed sample data

3. **Access the application**
   - **Web App**: http://localhost:3000
   - **API**: http://localhost:8080
   - **API Docs**: http://localhost:8080/swagger/
   - **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin123)

## 🛠️ Development

### Available Commands

```bash
# Development
task dev          # Start full dev environment
task up           # Start containers only
task down         # Stop containers
task logs         # View all logs

# Code Generation
task gen          # Generate all code from OpenAPI
task gen-go       # Generate Go server code
task gen-ts       # Generate TypeScript client

# Database
task migrate      # Run migrations
task migrate-down # Rollback migration
task seed         # Seed sample data

# Testing & Linting
task test         # Run all tests
task lint         # Run all linters
task api-test     # API tests only
task web-test     # Web tests only

# Building
task build        # Build all applications
task api-build    # Build API binary
task web-build    # Build Next.js app

# Health Check
task health       # Check all services
```

### Project Structure

```
GoGym/
├── apps/
│   ├── api/                     # Go API (Clean Architecture)
│   │   ├── cmd/api/            # Application entry point
│   │   ├── internal/
│   │   │   ├── domain/         # Domain entities & business logic
│   │   │   ├── usecase/        # Application services
│   │   │   └── adapter/        # Infrastructure adapters
│   │   ├── infra/
│   │   │   ├── migrations/     # Database migrations
│   │   │   └── seeds/          # Sample data
│   │   └── Dockerfile
│   └── web/                    # Next.js Frontend
│       ├── app/                # App Router pages
│       ├── components/         # React components
│       ├── lib/                # Utilities
│       ├── types/              # Generated TypeScript types
│       └── Dockerfile
├── packages/
│   └── openapi/                # OpenAPI specification
│       ├── openapi.yaml        # API contract
│       ├── gen-go.sh          # Go code generation
│       └── gen-ts.sh          # TypeScript code generation
├── infra/
│   └── docker/                 # Docker Compose configuration
├── gogym-old/                  # Original Rails application (backup)
└── ADRS/                       # Architecture Decision Records
```

## 🏛️ Architecture Principles

### Backend (Go + Clean Architecture)

- **Domain Layer**: Business entities, value objects, domain services
- **Use Case Layer**: Application services, input/output ports
- **Adapter Layer**: Infrastructure implementations (HTTP, database, external APIs)

### Frontend (Next.js + Server Components)

- **Server Components**: Data fetching, initial rendering
- **Client Components**: Interactive UI, user input handling
- **Server Actions**: Form submissions, mutations

### API-First Development

- OpenAPI 3.0 specification as single source of truth
- Generated Go server stubs and TypeScript client
- Contract testing and documentation

## 🗃️ Database Schema

### Key Tables

- **users**: User accounts with authentication
- **gyms**: Gym locations with spatial data (POINT SRID 4326)
- **reviews**: User reviews with ratings and photos
- **tags**: Categorization system for gyms
- **favorites**: User bookmarked gyms

### Spatial Features

- MySQL 8.0 with spatial indexing
- Location-based search using `ST_Distance_Sphere()`
- Full-text search on gym names and descriptions

## 🔧 Configuration

### Environment Variables

Key configuration in `.env`:

- **Database**: MySQL connection settings
- **Authentication**: JWT secrets and expiration
- **Storage**: S3/MinIO configuration for photos
- **External APIs**: Google Maps API key

### Feature Flags

- GraphQL endpoint (experimental)
- Social login providers
- Image optimization settings

## 🧪 Testing

### API Tests
```bash
cd apps/api
go test -v ./...
```

### Frontend Tests
```bash
cd apps/web
pnpm test
```

### Integration Tests
```bash
task up      # Start services
task test    # Run full test suite
```

## 📦 Deployment

### Development
- Docker Compose with local services
- Hot reloading for both API and Web

### Production (Recommended)
- **Frontend**: Vercel or similar Edge platform
- **Backend**: AWS App Runner, Google Cloud Run, or Kubernetes
- **Database**: AWS RDS Aurora MySQL, Google Cloud SQL
- **Storage**: AWS S3, Google Cloud Storage
- **Cache**: AWS ElastiCache, Google Memorystore

## 🤝 Contributing

1. Create feature branch from `main`
2. Make changes following our coding conventions
3. Run tests and linting: `task test lint`
4. Update OpenAPI spec if adding/changing APIs
5. Generate code: `task gen`
6. Submit pull request

### Code Standards

- **Go**: `gofumpt`, `golangci-lint`
- **TypeScript**: ESLint Flat Config, Biome/Prettier
- **Commits**: Conventional Commits format

## 📄 Documentation

- [ADRs](./ADRS/) - Architecture decisions
- [API Documentation](./packages/openapi/openapi.yaml) - OpenAPI spec
- [Database Schema](./apps/api/infra/migrations/) - Migration files
- [Original Rails App](./gogym-old/) - Previous implementation (backup)

## 🆘 Troubleshooting

### Common Issues

1. **Port conflicts**: Change ports in `docker-compose.yml`
2. **Permission errors**: Check Docker permissions
3. **Database connection**: Verify MySQL is running and accessible
4. **Code generation fails**: Check OpenAPI spec syntax

### Getting Help

- Check existing [Issues](../../issues)
- Review [ADRs](./ADRS/) for architectural context
- Run `task health` to check service status

---

Built with ❤️ for the fitness community in Japan