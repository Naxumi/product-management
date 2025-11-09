# Product Management System

[![Go](https://img.shields.io/badge/Go-1.25.4-00ADD8?logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-16.0.1-000000?logo=next.js)](https://nextjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?logo=postgresql)](https://www.postgresql.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A modern, full-stack product management application built with Go (backend) and Next.js (frontend), featuring comprehensive CRUD operations, image upload, advanced filtering, and multi-layer validation.

## ✨ Key Features

### Backend (Go)
- ✅ **Clean Architecture**: Domain-driven design with clear separation of concerns
- ✅ **RESTful API**: Chi router v5 with structured JSON responses
- ✅ **Comprehensive Testing**: 54 tests with 100% pass rate
  - 13 repository integration tests
  - 19 service unit tests with mocks
  - 22 HTTP handler tests
- ✅ **Advanced Validation**: Multi-layer validation with field-specific error messages
- ✅ **Image Management**: Local storage with UUID-based filenames, format & size validation
- ✅ **Precise Price Handling**: Decimal.Decimal for accurate financial calculations
- ✅ **Smart Filtering**: Name, SKU, category, status, and price range filters
- ✅ **Pagination**: Configurable page size with full metadata
- ✅ **CORS Support**: Pre-configured for frontend integration
- ✅ **Structured Logging**: Request/response middleware with JSON logging

### Frontend (Next.js)
- ✅ **Modern Stack**: Next.js 16 with App Router, React 19, TypeScript 5.x
- ✅ **Responsive Design**: Mobile-first UI with shadcn/ui components
- ✅ **Smart Error Handling**: ValidationError class for field-specific validation errors
- ✅ **Real-time Validation**: Frontend + backend validation with detailed feedback
- ✅ **Image Upload**: Drag-and-drop with preview, format validation (JPG/PNG/GIF), 5MB limit
- ✅ **Advanced Filtering**: Real-time search, category dropdown, status filter
- ✅ **Rich UI Components**: Dialogs, sheets, toasts, badges, dropdown menus
- ✅ **Indonesian Localization**: Toast notifications and error messages in Bahasa Indonesia
- ✅ **Optimized Performance**: Turbopack bundler with hot reload

## 🏗️ Architecture

### Backend Architecture (Go + PostgreSQL)

**Framework & Tools:**
- **Router**: Chi v5 - Lightweight, idiomatic HTTP router
- **Database**: PostgreSQL 14+ with pgx/v5 driver
- **Architecture Pattern**: Clean Architecture (Domain → Service → Repository → Handler)
- **Decimal Handling**: shopspring/decimal for precise price calculations
- **Testing**: testify/assert + mock for comprehensive test coverage

**Key Features:**
- Full CRUD operations with validation
- Image upload with local storage and multi-layer validation:
  - Format validation (JPG, JPEG, PNG, GIF)
  - Size validation (max 5MB)
  - MIME type verification
  - UUID-based unique filenames
- Advanced filtering and pagination
- Domain-specific error handling with custom error types
- CORS configuration for seamless frontend integration
- Static file serving for uploaded images
- Structured JSON logging with request/response middleware
- Transaction support for data integrity

**Error Handling:**
- `ValidationError`: Field-specific validation errors with details map
- Domain errors: `ErrProductNotFound`, `ErrProductSKUExists`, etc.
- Image errors: `ErrInvalidImageFormat`, `ErrImageTooLarge`, `ErrImageRequired`
- HTTP status mapping: 200, 201, 400, 404, 422, 500

### Frontend Architecture (Next.js + TypeScript)

**Framework & Libraries:**
- **Framework**: Next.js 16.0.1 with App Router and Turbopack
- **Language**: TypeScript 5.x for type safety
- **UI Components**: shadcn/ui (Dialog, Sheet, Table, Form, etc.)
- **Styling**: Tailwind CSS with custom theme
- **State Management**: React 19 hooks (useState, useEffect, useCallback)
- **Notifications**: Sonner for toast notifications
- **Form Handling**: Custom validation with API integration

**Key Features:**
- Product listing with responsive table layout
- Advanced filtering:
  - Real-time name search
  - Real-time SKU search  
  - Category dropdown filter
  - Status filter (All/Active/Inactive)
  - Reset filters functionality
- CRUD operations with validation:
  - **Create/Edit**: Dialog-based forms with inline validation
  - **Delete**: Confirmation dialog with product details
  - **View**: Sheet component for detailed product view
- Image management:
  - Upload with drag-and-drop or file picker
  - Real-time preview before upload
  - Format validation (JPG, JPEG, PNG, GIF)
  - Size validation (max 5MB)
  - Remove/delete image capability
  - Fullscreen image preview dialog
- Pagination with customizable page size (10, 25, 50, 100)
- **Smart Error Handling**:
  - `ValidationError` class for field-specific errors
  - Displays validation errors with field names
  - Toast notifications in Indonesian
  - Dismissible error alerts
- Responsive design with mobile navigation
- Loading states for all async operations

**Validation Flow:**
1. **Frontend Validation**: File picker restrictions + pre-submit validation
2. **API Client**: Detects `VALIDATION_ERROR` response code
3. **ValidationError**: Throws custom error with field details map
4. **Component Handling**: Displays field-specific error messages in toast
5. **Backend Validation**: Final security validation with domain errors

## 🚀 Quick Start

### Prerequisites
- **Go**: 1.21 or higher
- **Node.js**: 18 or higher  
- **PostgreSQL**: 14 or higher
- **golang-migrate**: For database migrations (optional but recommended)

### Backend Setup

1. **Navigate to backend directory:**
```bash
cd backend
```

2. **Create database:**
```sql
CREATE DATABASE product_management;
```

3. **Configure database connection** (if needed):
```go
// internal/config/config.go
// Default connection: postgres:root@localhost:5432/product_management
```

Or set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=root
export DB_NAME=product_management
```

4. **Run migrations:**
```bash
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management?sslmode=disable" up
```

5. **Run the server:**
```bash
go run cmd/api/main.go
```

✅ Server starts on `http://localhost:8080`

6. **Run tests (optional):**
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test suite
go test -v ./internal/service/product
```

### Frontend Setup

1. **Navigate to frontend directory:**
```bash
cd frontend
```

2. **Install dependencies:**
```bash
npm install
```

3. **Configure API URL** (optional):
```bash
# Create .env.local file
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local
```

4. **Run development server:**
```bash
npm run dev
```

✅ Frontend starts on `http://localhost:3000`

5. **Build for production (optional):**
```bash
npm run build
npm start
```

## 📡 API Reference

**Base URL:** `http://localhost:8080/api/v1`

### Endpoints

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| GET | `/product` | List products with filters & pagination | - | `ListProductResponse` |
| GET | `/product/{id}` | Get product by ID | - | `Product` |
| GET | `/product/sku/{sku}` | Get product by SKU | - | `Product` |
| POST | `/product` | Create new product | `CreateProductRequest` | `Product` |
| PUT | `/product` | Update existing product | `UpdateProductRequest` | `Product` |
| DELETE | `/product/{id}` | Delete product by ID | - | `Success` |
| POST | `/product/{id}/image` | Upload product image | `multipart/form-data` | `Success` |
| DELETE | `/product/{id}/image` | Delete product image | - | `Success` |
| GET | `/uploads/*` | Serve uploaded static files | - | Image file |

### Query Parameters (List Products)

| Parameter | Type | Description | Default | Example |
|-----------|------|-------------|---------|---------|
| `name` | string | Filter by name (partial match) | - | `?name=laptop` |
| `sku` | string | Filter by SKU (partial match) | - | `?sku=PROD` |
| `category` | string | Filter by category (exact match) | - | `?category=Electronics` |
| `status` | string | Filter by status | - | `?status=Active` |
| `min_price` | number | Minimum price filter | - | `?min_price=100000` |
| `max_price` | number | Maximum price filter | - | `?max_price=500000` |
| `page` | number | Page number | 1 | `?page=2` |
| `limit` | number | Items per page (max 100) | 10 | `?limit=25` |
| `sort_by` | string | Sort field | `created_at` | `?sort_by=price` |
| `sort_order` | string | Sort order (asc/desc) | `desc` | `?sort_order=asc` |

### Request/Response Examples

**Create Product:**
```json
POST /api/v1/product
{
  "sku": "PROD-001",
  "name": "Laptop Gaming",
  "description": "High-performance gaming laptop",
  "price": 15000000,
  "stock": 10,
  "category": "Electronics",
  "status": "Active"
}
```

**Success Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "sku": "PROD-001",
    "name": "Laptop Gaming",
    "description": "High-performance gaming laptop",
    "price": 15000000,
    "stock": 10,
    "category": "Electronics",
    "status": "Active",
    "image_url": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Validation Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "sku": "sku must be at least 3 characters",
      "price": "price must be greater than 0"
    }
  }
}
```

**Upload Image:**
```bash
POST /api/v1/product/1/image
Content-Type: multipart/form-data

image: [binary file data]
```

## 🧪 Testing

### Backend Test Suite

The backend includes **54 comprehensive tests** with **100% pass rate**, covering all layers of the application.

**Run all tests:**
```bash
cd backend
go test ./...
```

**Run with verbose output:**
```bash
go test -v ./...
```

**Run specific test suite:**
```bash
# Repository tests (integration with real database)
go test -v ./internal/repository/postgresql

# Service tests (unit tests with mocks)
go test -v ./internal/service/product

# Handler tests (HTTP endpoint tests)
go test -v ./internal/handler/http
```

**Test Coverage Breakdown:**
- ✅ **13 Repository Tests**: Integration tests with PostgreSQL
  - CRUD operations
  - Filter and pagination
  - Transaction handling
  - Constraint validation

- ✅ **19 Service Tests**: Unit tests with mocked dependencies
  - Business logic validation
  - Error handling
  - Image upload logic
  - Repository error propagation

- ✅ **22 Handler Tests**: HTTP endpoint tests
  - Request/response validation
  - Status code verification
  - JSON serialization
  - Error response format

**Sample Test Output:**
```bash
$ go test -v ./...

=== RUN   TestProductRepository_Create
--- PASS: TestProductRepository_Create (0.05s)
=== RUN   TestProductService_Create
--- PASS: TestProductService_Create (0.01s)
=== RUN   TestProductHandler_Create
--- PASS: TestProductHandler_Create (0.02s)

PASS
ok      github.com/naxumi/bnsp-jwd/internal/repository/postgresql    0.243s
ok      github.com/naxumi/bnsp-jwd/internal/service/product          0.087s
ok      github.com/naxumi/bnsp-jwd/internal/handler/http            0.156s
```

### Frontend Testing

Frontend testing setup (to be implemented):
```bash
cd frontend
npm test
```

## 📁 Project Structure

```
.
├── backend/                                # Go backend application
│   ├── cmd/
│   │   └── api/
│   │       └── main.go                     # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go                   # Configuration management
│   │   ├── domain/product/                 # Domain layer (business entities)
│   │   │   ├── entity.go                   # Product entity definition
│   │   │   ├── dto.go                      # Data transfer objects with validation
│   │   │   ├── errors.go                   # Domain-specific errors
│   │   │   ├── repository.go               # Repository interface
│   │   │   └── service.go                  # Service interface
│   │   ├── repository/postgresql/          # Data access layer
│   │   │   ├── product.go                  # Product repository implementation
│   │   │   ├── product_test.go             # Integration tests (13 tests)
│   │   │   └── transaction.go              # Transaction helper
│   │   ├── service/                        # Business logic layer
│   │   │   ├── product/
│   │   │   │   ├── service.go              # Product service implementation
│   │   │   │   └── service_test.go         # Unit tests with mocks (19 tests)
│   │   │   └── file/
│   │   │       └── service.go              # File upload service
│   │   ├── handler/http/                   # Presentation layer (HTTP handlers)
│   │   │   ├── product.go                  # Product HTTP handler
│   │   │   ├── product_test.go             # HTTP endpoint tests (22 tests)
│   │   │   ├── router.go                   # Route definitions
│   │   │   └── response/
│   │   │       ├── response.go             # Standard response helpers
│   │   │       └── error.go                # Error response mapping
│   │   └── pkg/                            # Shared packages
│   │       ├── database/
│   │       │   └── postgresql.go           # Database connection
│   │       ├── storage/
│   │       │   ├── storage.go              # Storage interface
│   │       │   └── local.go                # Local file storage implementation
│   │       └── validator/
│   │           └── validator.go            # Validation utilities
│   ├── migrations/                         # Database migrations
│   │   ├── 000001_create_products_table.up.sql
│   │   └── 000001_create_products_table.down.sql
│   ├── storage/                            # Uploaded files directory
│   │   └── products/                       # Product images organized by ID
│   └── go.mod                              # Go module dependencies
│
└── frontend/                               # Next.js frontend application
    ├── app/                                # Next.js App Router
    │   ├── layout.tsx                      # Root layout with Header/Footer
    │   ├── globals.css                     # Global styles and Tailwind config
    │   ├── page.tsx                        # Home/landing page
    │   └── products/
    │       └── page.tsx                    # Products CRUD page (main application)
    ├── components/                         # React components
    │   ├── ui/                             # shadcn/ui components
    │   │   ├── alert.tsx                   # Alert/notification component
    │   │   ├── badge.tsx                   # Status badge component
    │   │   ├── button.tsx                  # Button component
    │   │   ├── card.tsx                    # Card container
    │   │   ├── dialog.tsx                  # Modal dialog
    │   │   ├── dropdown-menu.tsx           # Dropdown action menu
    │   │   ├── form.tsx                    # Form components
    │   │   ├── input.tsx                   # Text input field
    │   │   ├── label.tsx                   # Form label
    │   │   ├── pagination.tsx              # Pagination controls
    │   │   ├── select.tsx                  # Dropdown select
    │   │   ├── sheet.tsx                   # Slide-out panel
    │   │   ├── sonner.tsx                  # Toast notification wrapper
    │   │   ├── table.tsx                   # Table component
    │   │   └── textarea.tsx                # Multi-line text input
    │   ├── header.tsx                      # Application header with logo
    │   ├── footer.tsx                      # Application footer with links
    │   ├── product-form-dialog.tsx         # Create/Edit product dialog
    │   ├── delete-confirm-dialog.tsx       # Delete confirmation dialog
    │   ├── image-upload.tsx                # Image upload with validation
    │   └── image-preview-dialog.tsx        # Fullscreen image preview
    ├── lib/                                # Utilities and helpers
    │   ├── api.ts                          # API client with all endpoints
    │   ├── types.ts                        # TypeScript interfaces & ValidationError class
    │   └── utils.ts                        # Utility functions (cn, etc.)
    ├── public/                             # Static assets
    │   └── logo.png                        # Company logo
    ├── .env.local                          # Environment variables (not in git)
    ├── components.json                     # shadcn/ui configuration
    ├── tailwind.config.ts                  # Tailwind CSS configuration
    ├── tsconfig.json                       # TypeScript configuration
    ├── next.config.ts                      # Next.js configuration
    └── package.json                        # Node.js dependencies
```

## 🎯 Features in Detail

### 1. Smart Error Handling & Validation

**Multi-Layer Validation Architecture:**

```
User Input → Frontend Validation → API Request → Backend Validation → Database
     ↓              ↓                    ↓              ↓                ↓
  Field Rules   Format Check      ValidationError   Domain Logic    Constraints
```

**Backend (Go):**
- `ValidationError` struct with field-specific error messages
- `ValidationErrors` slice that converts to map for JSON responses
- Domain errors: `ErrProductNotFound`, `ErrProductSKUExists`, etc.
- HTTP status code mapping (422 for validation, 404 for not found, etc.)

```go
// Example: Field-specific validation error
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "sku": "sku must be at least 3 characters",
      "price": "price must be greater than 0"
    }
  }
}
```

**Frontend (TypeScript):**
- `ValidationError` class that extends Error
- API client detects `VALIDATION_ERROR` response code
- Throws ValidationError with field details map
- Components display field-specific errors in toast notifications

```typescript
// Example: ValidationError handling
try {
  await apiClient.createProduct(data);
} catch (err) {
  if (err instanceof ValidationError) {
    // Display: "sku: sku must be at least 3 characters, price: price must be greater than 0"
    const fieldErrors = Object.entries(err.details)
      .map(([field, message]) => `${field}: ${message}`)
      .join(", ");
    toast.error("Validasi Gagal", { description: fieldErrors });
  }
}
```

### 2. Image Upload & Management

**Multi-Layer Validation:**
1. **Frontend File Picker**: Restricts to `image/jpeg,image/png,image/gif`
2. **Frontend Pre-Upload**: Validates type, extension, and size before submission
3. **Backend MIME Type**: Verifies Content-Type header
4. **Backend Extension**: Case-insensitive check (.jpg, .jpeg, .png, .gif)
5. **Backend Size**: Enforces 5MB limit

**Storage Structure:**
```
storage/products/{productId}/{uuid}-{originalFilename}.{ext}
```

**Features:**
- UUID-based filenames to prevent conflicts
- Automatic directory creation per product
- Old image deletion on new upload
- Public URL generation for frontend access
- Domain-specific errors: `ErrInvalidImageFormat`, `ErrImageTooLarge`, `ErrImageRequired`

### 3. Advanced Filtering & Pagination

**Filter Options:**
- **Text Search**: Name and SKU with partial matching (case-insensitive)
- **Category Filter**: Exact match dropdown selection
- **Status Filter**: All, Active, or Inactive
- **Price Range**: Min and max price filters with validation
- **Reset**: Clear all filters with one click

**Pagination Features:**
- Configurable page size: 10, 25, 50, or 100 items
- Page navigation with first/previous/next/last buttons
- Total count and page info display
- Query string persistence for bookmarking

**Backend Query Optimization:**
- Dynamic WHERE clause building
- Parameterized queries to prevent SQL injection
- Efficient COUNT query for total records
- Sort by any field with ASC/DESC order

### 4. Precise Price Handling

**Decimal Precision:**
- Backend uses `shopspring/decimal` library for financial calculations
- Database stores as `NUMERIC(10,2)` for exact decimal representation
- Frontend formats as Indonesian Rupiah (IDR): `Rp 15.000.000`
- Max price validation: `Rp 99.999.999,99`

**Currency Formatting:**
```typescript
// Frontend: Auto-format input as user types
const formatIDR = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value);
};
```

### 5. Responsive Design & UX

**Mobile-First Approach:**
- Responsive table with horizontal scroll on mobile
- Mobile navigation with Sheet (slide-out) menu
- Touch-friendly button sizes (min 44x44px)
- Optimized image loading with lazy loading

**User Feedback:**
- Toast notifications for all actions (success/error)
- Loading states on all buttons during async operations
- Dismissible error alerts with X button
- Inline validation errors in forms
- Indonesian language for better local UX

### 6. Testing & Quality Assurance

**Backend Test Coverage: 54 Tests**

**Repository Layer (13 tests):**
- Create product with valid data
- Duplicate SKU constraint violation
- Get product by ID and SKU
- Update product fields
- Delete product
- List with filters (name, SKU, category, status, price range)
- Pagination and sorting

**Service Layer (19 tests):**
- Create with validation
- Get by ID and SKU
- Update business logic
- Delete with error handling
- Upload image with format/size validation
- Repository error propagation

**Handler Layer (22 tests):**
- HTTP status codes (200, 201, 400, 404, 422, 500)
- JSON request parsing
- JSON response serialization
- Validation error response format
- Image upload via multipart/form-data

### 7. Clean Architecture Benefits

**Separation of Concerns:**
- **Domain**: Business entities and rules (no external dependencies)
- **Service**: Business logic and orchestration
- **Repository**: Data access and persistence
- **Handler**: HTTP presentation and routing

**Dependency Flow:**
```
Handler → Service → Repository → Database
   ↓         ↓          ↓
Response   Domain    Entity
```

**Benefits:**
- ✅ Easy to test each layer independently
- ✅ Swap implementations (e.g., PostgreSQL → MySQL)
- ✅ Clear responsibility for each component
- ✅ Reduced coupling between layers
- ✅ Better code organization and maintainability

## ⚙️ Configuration

### Backend Environment Variables

Configure via environment variables or edit `internal/config/config.go`:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `root` |
| `DB_NAME` | Database name | `product_management` |
| `SERVER_PORT` | HTTP server port | `8080` |

**Example:**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=root
export DB_NAME=product_management
export SERVER_PORT=8080
```

### Frontend Environment Variables

Create `.env.local` in the frontend directory:

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API base URL | `http://localhost:8080/api/v1` |

**Example `.env.local`:**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### CORS Configuration

Backend CORS is pre-configured for local development:

```go
// internal/handler/http/router.go
cors.Handler(cors.Options{
    AllowedOrigins:   []string{"http://localhost:3000"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Content-Type"},
    AllowCredentials: false,
    MaxAge:           300,
})
```

For production, update `AllowedOrigins` to your frontend domain.

### Database Schema

The product table schema:

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    category VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('Active', 'Inactive')),
    image_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_created_at ON products(created_at DESC);
```

## 🚀 Production Deployment

### Backend Deployment

1. **Build the application:**
```bash
cd backend
go build -o product-api cmd/api/main.go
```

2. **Set environment variables:**
```bash
export DB_HOST=your-db-host
export DB_PORT=5432
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=product_management
export SERVER_PORT=8080
```

3. **Run the binary:**
```bash
./product-api
```

4. **Using systemd (Linux):**
```ini
# /etc/systemd/system/product-api.service
[Unit]
Description=Product Management API
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/product-management/backend
Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=postgres"
Environment="DB_PASSWORD=your-password"
Environment="DB_NAME=product_management"
Environment="SERVER_PORT=8080"
ExecStart=/opt/product-management/backend/product-api
Restart=always

[Install]
WantedBy=multi-user.target
```

### Frontend Deployment

#### Option 1: Vercel (Recommended)

1. Push code to GitHub
2. Import repository in Vercel
3. Set environment variable:
   - `NEXT_PUBLIC_API_URL`: Your backend API URL
4. Deploy automatically

#### Option 2: Self-Hosted

1. **Build the application:**
```bash
cd frontend
npm run build
```

2. **Start production server:**
```bash
npm start
```

3. **Using PM2:**
```bash
npm install -g pm2
pm2 start npm --name "product-frontend" -- start
pm2 save
pm2 startup
```

#### Option 3: Static Export

1. **Update `next.config.ts`:**
```typescript
const nextConfig = {
  output: 'export',
};
```

2. **Build static files:**
```bash
npm run build
```

3. **Deploy `out/` folder to any static hosting (Netlify, GitHub Pages, etc.)**

### Docker Deployment (Optional)

**Backend Dockerfile:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o product-api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/product-api .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/storage ./storage
EXPOSE 8080
CMD ["./product-api"]
```

**Frontend Dockerfile:**
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
EXPOSE 3000
CMD ["npm", "start"]
```

**Docker Compose:**
```yaml
version: '3.8'

services:
  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: root
      POSTGRES_DB: product_management
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: root
      DB_NAME: product_management
      SERVER_PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - db

  frontend:
    build: ./frontend
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080/api/v1
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  postgres_data:
```

## 🐛 Troubleshooting

### Database Connection Issues

**Problem**: Cannot connect to PostgreSQL

**Solutions:**
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql  # Linux
brew services list                # macOS
Get-Service postgresql*           # Windows (PowerShell)

# Test connection
psql -U postgres -h localhost -d product_management

# Verify credentials in config
# backend/internal/config/config.go
```

### CORS Errors

**Problem**: CORS policy blocking frontend requests

**Solutions:**
```go
// Update backend/internal/handler/http/router.go
cors.Handler(cors.Options{
    AllowedOrigins: []string{
        "http://localhost:3000",
        "https://your-production-domain.com",
    },
    // ...
})
```

### Image Upload Failures

**Problem**: Image upload returns error or fails silently

**Common Causes & Solutions:**

1. **Directory permissions:**
```bash
# Ensure storage directory is writable
mkdir -p backend/storage/products
chmod 755 backend/storage/products
```

2. **File size too large:**
   - Check file is < 5MB
   - Frontend validates before upload
   - Backend also validates

3. **Invalid file format:**
   - Only JPG, JPEG, PNG, GIF allowed (case-insensitive)
   - Check MIME type matches extension

4. **Network timeout:**
   - Large files may timeout on slow connections
   - Consider increasing timeout in HTTP client

### Tests Failing

**Problem**: Backend tests fail

**Solutions:**
```bash
# Ensure test database is accessible
createdb product_management_test

# Update test database config if needed
# Check internal/repository/postgresql/product_test.go

# Run migrations on test database
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management_test?sslmode=disable" up

# Clean dependencies
go mod tidy
go clean -testcache

# Run tests with verbose output
go test -v ./...
```

### Frontend Build Errors

**Problem**: Build fails or shows TypeScript errors

**Solutions:**
```bash
# Clear Next.js cache
rm -rf .next

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install

# Check TypeScript errors
npm run build

# Update environment variables
cp .env.local.example .env.local
```

### Port Already in Use

**Problem**: Port 8080 or 3000 already in use

**Solutions:**
```bash
# Find process using port (Linux/macOS)
lsof -i :8080
lsof -i :3000

# Find process using port (Windows)
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# Kill process
kill -9 <PID>           # Linux/macOS
taskkill /PID <PID> /F  # Windows

# Or change port in config
# Backend: SERVER_PORT environment variable
# Frontend: npm run dev -- -p 3001
```

### Migration Errors

**Problem**: Migration fails or database schema mismatch

**Solutions:**
```bash
# Check migration status
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management?sslmode=disable" version

# Rollback migration
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management?sslmode=disable" down

# Reapply migrations
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management?sslmode=disable" up

# Force version (use with caution)
migrate -path migrations -database "postgresql://postgres:root@localhost:5432/product_management?sslmode=disable" force <version>
```

## 🛣️ Development Roadmap

### Planned Features

- [ ] **Authentication & Authorization**
  - User registration and login
  - JWT-based authentication
  - Role-based access control (Admin, User)
  - Protected routes

- [ ] **Product Categories Management**
  - CRUD operations for categories
  - Category hierarchy/nesting
  - Category-based filtering

- [ ] **Inventory Tracking**
  - Stock movement history
  - Low stock alerts
  - Automatic reorder points

- [ ] **Advanced Search**
  - Full-text search
  - Elasticsearch integration
  - Search suggestions/autocomplete

- [ ] **Bulk Operations**
  - Bulk product import (CSV/Excel)
  - Bulk product export
  - Bulk update/delete

- [ ] **Product Variants**
  - Size, color, material variants
  - Variant-specific pricing
  - Variant-specific stock

- [ ] **Analytics Dashboard**
  - Sales trends
  - Top products
  - Stock levels overview
  - Real-time metrics

- [ ] **Multi-Language Support**
  - English, Indonesian, and more
  - i18n integration
  - Language switcher

- [ ] **API Documentation**
  - Swagger/OpenAPI integration
  - Interactive API explorer
  - Code generation

- [ ] **Performance Optimization**
  - Redis caching
  - Database query optimization
  - Image CDN integration
  - Lazy loading improvements

### Completed Features

- ✅ Product CRUD operations
- ✅ Image upload with validation
- ✅ Advanced filtering and pagination
- ✅ Multi-layer validation with field-specific errors
- ✅ Clean Architecture implementation
- ✅ Comprehensive test coverage (54 tests)
- ✅ Responsive design
- ✅ Indonesian localization for UI messages

## 📚 Learn More

### Technologies Used

**Backend:**
- [Go](https://go.dev/) - Programming language
- [Chi Router](https://github.com/go-chi/chi) - HTTP router
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [shopspring/decimal](https://github.com/shopspring/decimal) - Decimal numbers
- [testify](https://github.com/stretchr/testify) - Testing toolkit

**Frontend:**
- [Next.js](https://nextjs.org/) - React framework
- [TypeScript](https://www.typescriptlang.org/) - Type safety
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [shadcn/ui](https://ui.shadcn.com/) - UI components
- [Sonner](https://sonner.emilkowal.ski/) - Toast notifications

**Database:**
- [PostgreSQL](https://www.postgresql.org/) - Relational database
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations

### Recommended Reading

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [Next.js Documentation](https://nextjs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Development Guidelines:**
- Write tests for new features
- Follow existing code style
- Update documentation
- Ensure all tests pass before submitting PR

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Muhammad Hafizh Al Furqon**

- GitHub: [@Naxumi](https://github.com/Naxumi)
- Email: muhammad.hafizhalfurqon@gmail.com

## 🙏 Acknowledgments

- Thanks to the Go community for excellent libraries
- shadcn for beautiful UI components
- Vercel for Next.js and hosting platform
- PostgreSQL team for robust database system

---

**Built with ❤️ using Go, Next.js, and PostgreSQL**

*Last updated: November 2024*
