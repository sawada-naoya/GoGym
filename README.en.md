# GoGym

GoGym is a web application designed for managing workout records.

---

## Key Features

- Workout record management (body parts / exercises / weight / reps)
- Monthly view for workout history
- Mobile-optimized UI
- Authentication
- Multi-language support (Japanese / English)

---

## Design Characteristics

- Separated architecture with Go API server and Next.js web frontend
- Go backend designed based on Clean Architecture
- Update flow using Next.js App Router + Server Actions
- Frontend structure with responsibilities separated by features
- Prevention of BFF bloat through actions / apis separation

---

## Tech Stack

### Frontend

- Next.js 15 (App Router)
- React
- TypeScript
- Tailwind CSS
- NextAuth v5
- Server Actions

### Backend

- Go 1.25
- Echo (Web Framework)
- Clean Architecture
- Air (Hot reload)

### Infrastructure

- PostgreSQL 16
- Docker / Docker Compose

---

## Development Environment Setup

### Prerequisites

- Docker & Docker Compose
- Node.js 20 or higher
- Go 1.25
- Air (Go hot reload tool)

### Getting Started

#### 1. Start PostgreSQL

```bash
# Clone the repository
git clone <repository-url>
cd GoGym

# Copy environment variables file
cd infra
cp .env.sample .env

# Start PostgreSQL container
docker-compose up -d postgres
```

#### 2. Start API Server (separate terminal)

```bash
cd apps/api
air
```

#### 3. Start Web Frontend (separate terminal)

```bash
cd apps/web
npm install
npm run dev
```

### Access

- **Web**: http://localhost:3003
- **API**: http://localhost:8081
- **PostgreSQL**: localhost:5433

### Stopping the Services

```bash
# Stop PostgreSQL
cd infra
docker-compose down

# Stop API / Web with Ctrl + C in each terminal
```

---

## License

MIT
