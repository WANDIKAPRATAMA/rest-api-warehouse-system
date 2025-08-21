# Warehouse Management System

A RESTful API built with Go and Fiber for managing warehouse inventory, products, and user authentication.

## Overview

This application provides a modular service for managing warehouse locations, product categories, products, stock movements, and a dashboard with detailed analytics. It uses a PostgreSQL database with GORM as the ORM and follows a clean architecture with separation of concerns.

## Prerequisites

- Go (version 1.21+)
- PostgreSQL (version 15+)
- Git

## Database Setup

Create a database named `db_warehouse` and apply the following schema:

### Schema

```sql
CREATE DATABASE db_warehouse;

\c db_warehouse

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    status user_status NOT NULL DEFAULT 'inactive',
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_user_id UUID NOT NULL,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    avatar_url TEXT,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (source_user_id) REFERENCES users(id)
);

CREATE TABLE user_security (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_user_id UUID NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (source_user_id) REFERENCES users(id)
);

CREATE TABLE application_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_user_id UUID NOT NULL,
    role app_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (source_user_id) REFERENCES users(id)
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_user_id UUID NOT NULL,
    device_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP,
    revoked_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (source_user_id, device_id)
);

CREATE TYPE stock_status AS ENUM ('available', 'low-stock', 'out-of-stock');
CREATE TYPE movement_type AS ENUM ('inbound', 'outbound');
CREATE TYPE user_status AS ENUM ('active', 'inactive');
CREATE TYPE app_role AS ENUM ('user', 'admin', 'super_admin');

CREATE TABLE product_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    category_id UUID NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    FOREIGN KEY (category_id) REFERENCES product_categories(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE warehouse_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE product_stocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_product_id UUID NOT NULL,
    warehouse_location_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    status stock_status DEFAULT 'available',
    updated_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (source_product_id) REFERENCES products(id),
    FOREIGN KEY (warehouse_location_id) REFERENCES warehouse_locations(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_product_id UUID NOT NULL,
    movement_type movement_type NOT NULL,
    quantity INT NOT NULL,
    reference_note TEXT,
    created_by UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (source_product_id) REFERENCES products(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```
2. Set up the database:
   - Create `db_warehouse` database in PostgreSQL.
   - Apply the schema above using a SQL client or migration tool.
3. Configure environment variables:
   - Create a `.env` file with:
     ```
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=your_user
     DB_PASSWORD=your_password
     DB_NAME=db_warehouse
     ```
4. Run the application:
   ```bash
   cd cmd/web
   go run .
   ```
5. Access the API at `http://localhost:3000/api/...` (default port, adjust in code if needed).

## Next Steps

For detailed architecture and API documentation, refer to [ARCHITECTURE.md](ARCHITECTURE.md) and [ROUTES.md](ROUTES.md).

```

```
