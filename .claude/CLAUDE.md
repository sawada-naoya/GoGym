# CLAUDE.md

You must refactor this repository

## As Is

## Project Overview

GoGym is a gym search and review service built with Ruby on Rails 7.0.8 and React. It helps users find suitable gyms through location-based search, reviews, and recommendation features. The service targets people who want to start training or find better gym facilities.

## Technology Stack

- **Backend**: Ruby on Rails 7.0.8, Ruby 3.2.2
- **Frontend**: HTML, CSS, JavaScript, HotWire (Turbo/Stimulus), React 18.2.0
- **Database**: PostgreSQL
- **Infrastructure**: Docker, Render (production)
- **APIs**: Google Maps JavaScript API, Geocoding API
- **Key Gems**: Sorcery (authentication), Ransack (search), Kaminari (pagination), Draper (decorators), CarrierWave (file uploads), Geocoder, Sidekiq (background jobs)

## Development Commands

### Rails/Ruby Commands

```bash
# Start the Rails server
bundle exec rails server

# Run database migrations
bundle exec rails db:migrate

# Reset database (development)
bundle exec rails db:drop db:create db:migrate db:seed

# Run tests
bundle exec rspec

# Run RuboCop (linting)
bundle exec rubocop

# Install gems
bundle install
```

### JavaScript/Frontend Commands

```bash
# Build JavaScript bundles
yarn build

# Build CSS
yarn build:css

# Watch CSS changes during development
yarn watch:css

# Install npm packages
yarn install
```

### Docker Commands

```bash
# Start all services (database + web)
docker-compose up

# Start in development mode
docker-compose -f docker-compose-dev.yml up

# Build and start
docker-compose up --build
```

## Architecture and Structure

### Core Models and Relationships

- **User**: Authentication via Sorcery, has many reviews and favorites
- **Gym**: Central entity with location data, has many reviews and tags
- **Review**: User reviews of gyms with ratings and photos
- **Location**: Geographic data for gyms using Geocoder
- **Tag/GymTag**: Tagging system for gym categorization
- **Favorite**: User bookmarking of preferred gyms

### Key Directories

- `app/controllers/` - Rails controllers handling HTTP requests
- `app/models/` - ActiveRecord models and business logic
- `app/decorators/` - Draper decorators for presentation logic
- `app/services/` - Service objects for complex business operations
- `app/javascript/` - Stimulus controllers and React components
- `app/views/` - ERB templates with HotWire integration
- `spec/` - RSpec test suite with factories, helpers, and request specs

### Frontend Architecture

- **HotWire**: Primary frontend framework using Turbo and Stimulus
- **React**: Used for specific interactive components
- **Bootstrap 5**: UI framework with custom SCSS
- **esbuild**: JavaScript bundling
- **Stimulus controllers**: Located in `app/javascript/controllers/`

### Key Features Implementation

- **Map-based search**: Google Maps API integration in views and JavaScript
- **Geocoding**: Server-side location processing using Geocoder gem
- **Image uploads**: CarrierWave with AWS S3 storage
- **Background jobs**: Sidekiq for asynchronous processing
- **Internationalization**: Rails i18n with Japanese (ja) and English (en) locales
- **Search and filtering**: Ransack gem for complex queries
- **Pagination**: Kaminari for result pagination

## Configuration Notes

- **Timezone**: Set to 'Asia/Tokyo'
- **Default locale**: Japanese (:ja)
- **Database**: PostgreSQL with timezone-aware settings
- **File uploads**: Configured for AWS S3 in production
- **Background processing**: Redis and Sidekiq configured
- **TypeScript**: Configured for React components in `app/javascript/`

## Testing

- **Framework**: RSpec with FactoryBot
- **Structure**: Request specs, model specs, decorator specs, helper specs
- **Test database**: Separate PostgreSQL database for testing
- **Coverage**: Tests located in `spec/` directory with corresponding structure to `app/`

This project follows Rails conventions with modern frontend tooling, emphasizing location-based features and user-generated content through reviews and ratings.

## To Be

You are a senior software engineer. Design and scaffold a modern web product from scratch.
Generate actual files with code (path-prefixed code blocks). Avoid descriptions-only.

## Product

- Name: GoGym Next (codename: `gogym-next`)
- Purpose: Gym search & review (map search, reviews, photos, tags, favorites)
- Region: Japan (Asia/Tokyo), i18n (ja/en)

## Tech Constraints (must)

- Frontend: Next.js 14/15 App Router + React Server Components, TypeScript 5.x, Server Actions (edge-safe where possible)
- Backend: Go 1.22+ with Echo or Chi, **Clean Architecture / Onion**
- DB: **MySQL 8.0**（既存テーブル名・カラム名は維持）
- Jobs: **Asynq** (Go + Redis)
- Storage: S3 compatible（dev: MinIO） with presigned PUT
- Contract-first: **OpenAPI 3.0**（oapi-codegen で Go server stubs + TS client）
  - Optional **GraphQL (gqlgen)** as an experimental, feature-flagged endpoint
- Observability: **OpenTelemetry traces** (FE→BFF→DB) + Sentry
- Auth: email/password (argon2id), JWT (access/refresh rotation), Cookie (Secure/HttpOnly/SameSite)
- Container: **Docker Compose** first; prod 想定は Vercel(Front) + App Runner/ECS(Back) + RDS (Aurora MySQL) + ElastiCache + S3 + CloudFront
- CI: GitHub Actions（lint/test/build/OpenAPI drift check）

## Modern Coding Conventions (enforce)

- **Go**
  - stdlog: `log/slog` 構造化ログ + request-scoped logger (trace/span IDs)
  - errors: `errors.Join/Is/As`, `%w` wrapping、ハンドラでは domain error → HTTP へ明示的マッピング
  - DB: GORM（既存知見）＋**sqlc 併用**（読み/レポート系の型安全・高速化）
  - config: env（12-factor）+ minimal `envconfig` / `sops`対応の余地
  - tests: `testify`, `httpexpect`, repository は DB コンテナに対する integration test も用意
  - lint/format: `gofumpt`, `gci`, `golangci-lint`
- **TypeScript/Next.js**
  - ESLint **Flat config**、`biome` で整形（or Prettier、どちらか 1 つに統一）
  - RSC/Server Actions 優先。クライアントは UI/interaction に限定
  - データ型は **OpenAPI 生成の TS client** を単一真実源に
  - コンポーネントは `use client` を最小化、フォームは Server Action + progressive enhancement
  - e2e: Playwright、単体: Vitest + React Testing Library
  - i18n: `next-intl`（RSC 対応）
- **Monorepo**
  - `pnpm` + optional Turborepo（build 並列化）
  - `Taskfile.yml` or `Makefile` で、`gen`, `up`, `migrate`, `test` をワンコマンド化
- **Docs**
  - ADR（0001-openapi-first、0002-db-mysql-geo、0003-graphql-experiment）

## Architecture (Clean/Onion boundaries)

1. **Repo layout** (monorepo)

```
GoGym/
  apps/
    api/                      # Go (Echo/Chi). Clean/Onion 構成
      cmd/api/                # Composition root (DI/boot)
      internal/
        domain/               # エンティティ/値オブジェクト/ドメインサービス/エラー
          gym/
          review/
          user/
          common/
        usecase/              # アプリケーションサービス（入力/出力ポート）
          gym/
          review/
          user/
        adapter/              # インフラ適用層（外界への適合）
          http/               # ハンドラ（REST/GraphQLResolver）
          db/                 # GORM/sqlc 実装
            gorm/
            sqlc/
          job/                # Asynq タスク定義
          auth/               # JWT/Session 発行・検証
        gen/                  # oapi-codegen / gqlgen の生成物（コミット対象）
      infra/
        migrations/           # goose/atlas の DDL
        seeds/                # 開発用 seed
      configs/                # api.yaml/otel.yaml など
      Dockerfile
      Makefile                # api 単体向け（ルート Taskfile からも叩く）
    web/                      # Next.js (App Router, RSC)
      app/
        (search)/             # 検索ページ(RSC)
        api/                  # Route handlers（必要最小限）
        layout.tsx
        page.tsx
      components/
      lib/
      types/                  # OpenAPI から生成された TS 型のみ（手書きは置かない）
      public/
      next.config.ts
      package.json
      Dockerfile
  packages/
    openapi/
      openapi.yaml            # 契約の単一の真実源
      gen-go.sh               # oapi-codegen: Go（types/server/client）
      gen-ts.sh               # openapi-typescript: TS クライアント
      README.md
  infra/
    docker/
      docker-compose.yml      # mysql, redis, minio, imgproxy, api, web
      .env.sample
    terraform/                # 後日: RDS/S3/CF/WAF（任意）
  ADRS/
    0001-openapi-first.md
    0002-db-mysql-geo.md
    0003-graphql-experiment.md
  .github/
    workflows/
      ci.yml                  # lint/test/build/openapi drift
      preview.yml             # プレビュー/タグ配信（任意）
  Taskfile.yml                # モノレポ操作（gen/up/migrate/test 等）
  README.md

```

## Database: **MySQL 8.0**（カラム名・テーブル名は既存を踏襲）

- `gyms(id, name, description, location POINT SRID 4326, address, city, prefecture, postal_code, fts TEXT GENERATED, created_at, updated_at)`
  - `fts` = `CONCAT(COALESCE(name,''),' ',COALESCE(description,''))`（VIRTUAL 生成列）
  - Index: `SPATIAL INDEX(location)`, `FULLTEXT(fts)`
  - 近傍検索: `ST_Distance_Sphere(location, ST_SRID(POINT(:lon,:lat),4326)) <= :radius_m`
- `users(id, email VARCHAR(255) UNIQUE, password_hash, display_name, created_at, updated_at)`（email は `utf8mb4_0900_ai_ci` など CI コレーション）
- `reviews(id, user_id, gym_id, rating 1..5, comment, photos JSON, created_at, updated_at)`（FK index: user_id, gym_id）
- `tags(id, name UNIQUE)`, `gym_tags(gym_id, tag_id)`（PK: (gym_id, tag_id)）
- Migration: **atlas** または **goose**（どちらでも良いが自動化スクリプト必須）
- Seed: 最小サンプル（3 gyms, 2 users, 3 reviews, tags）

## APIs (OpenAPI-first; minimal)

- `GET /gyms/search?q&lat&lon&radius_m&cursor&limit` … **keyset pagination**
- `GET /gyms/{id}`
- `POST /auth/signup`
- `POST /auth/login`
- `POST /photos/presign`（S3 presigned PUT）
  → oapi-codegen で Go server stubs + TS client を生成（**単一真実源の型**）

## GraphQL (experiment; feature-flagged)

- gqlgen で `/graphql` を別パス or 別ポートに
- Schema (minimal):
  - Query: `gyms(q, lat, lon, radiusM, cursor, limit)`, `gym(id)`
  - Mutation: `createReview(gymId, rating, comment, photos)`
- DI/Usecase を共通にし、Resolver は handler ラッパに徹する

## Deliverables (generate **actual files** now)

1. **Repo skeleton**（上記ツリー）
2. `packages/openapi/openapi.yaml`（最小仕様、components/schemas も定義）
3. Codegen scripts:
   - Go: `oapi-codegen`（types, chi-server or echo-server, client）
   - TS: `openapi-typescript` or `oapi-codegen` TS client
4. **Docker Compose**（api, web, mysql:8.0, redis:7, minio, imgproxy）
5. **Migrations**（MySQL 8.0; POINT SRID 4326, 生成列 fts, indexes）
6. Go API:
   - `/gyms/search` ハンドラ（RSC で叩ける最低限; sqlc or GORM 実装例）
   - Auth skeleton（argon2id hash, JWT issue/verify, refresh rotation with Redis）
   - S3 presign endpoint（MinIO dev）
   - `slog` + `otel` 初期化
7. Next.js:
   - App Router 初期ページ（検索フォーム）+ **Server Action** で BFF 呼び出し
   - `next-intl` 設定（ja/en）
8. Quality:
   - **ESLint Flat** + biome (or prettier) 設定
   - Go: `golangci-lint` 設定、`Make/Task` ターゲット
   - GitHub Actions: lint/test/build、OpenAPI drift check（`git diff` on openapi）
9. ADR:
   - `ADRS/0001-openapi-first.md`
   - `ADRS/0002-db-mysql-geo.md`
   - `ADRS/0003-graphql-experiment.md`
10. **Makefile/Taskfile**:

- `gen`（OpenAPI TS/Go 双方向生成）
- `up` `down` `logs` `migrate` `api` `web` `test`

## Acceptance checklist

- `docker compose up` で api/web/db が起動、`GET /gyms/search` が 200 で返る
- OpenAPI から TS/Go の生成が通る（CI でも）
- `gyms.location` が SRID 4326 で保存され、`SPATIAL INDEX` が作成済
- Next.js (RSC) の Server Action から BFF を 1 回叩ける
- Lint & minimal test が通る
- ADR 3 本が存在する

> Output: Provide code blocks per file with `path:` headers and runnable minimal content.
