# GodocAPI

GodocAPI is a high-performance, Go-based RESTful API service designed for managing documents. It leverages the Fiber framework for HTTP handling, PostgreSQL for metadata storage, and RustFS (S3-compatible) for object storage.

## Features

-   **Document Management**: Upload, download, list, and delete documents.
-   **High Performance**: Built with Go 1.22 and Fiber v2.
-   **Metadata Storage**: Reliable metadata management using PostgreSQL.
-   **Object Storage**: Scalable file storage using RustFS (via AWS SDK for Go).
-   **Docker Ready**: Includes Dockerfile for easy containerization.

## Tech Stack

-   **Language**: Go 1.22
-   **Framework**: Fiber v2
-   **Database**: PostgreSQL (pgx driver)
-   **Storage**: RustFS (S3 Compatible)
-   **Config**: godotenv

## Prerequisites

-   Go 1.22+
-   PostgreSQL 14+
-   RustFS (or any S3-compatible storage like MinIO)

## Configuration

Duplicate the `.env` file or set environment variables:

```properties
SERVER_PORT=:8080
DB_URL=postgres://user:password@localhost:5432/godocapi?sslmode=disable
RUSTFS_ENDPOINT=http://localhost:9000
RUSTFS_ACCESS_KEY=your_access_key
RUSTFS_SECRET_KEY=your_secret_key
RUSTFS_BUCKET=documents
RUSTFS_REGION=us-east-1
```

## Database & Storage Setup

The easiest way to set up PostgreSQL and RustFS locally is using Docker Compose.

1.  **Start dependencies**:
    ```bash
    docker-compose up -d
    ```
    This will start:
    -   PostgreSQL on port `5432`
    -   RustFS (Object Storage) on port `9000` (Console: `9001`)
    -   It will also automatically create the `documents` bucket.

2.  **Initialize Database Schema**:
    Use a tool like DBeaver or `psql` to run the following DDL:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS documents (
  id           UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
  filename     TEXT        NOT NULL,
  storage_path TEXT        NOT NULL UNIQUE,
  size         BIGINT      NOT NULL CHECK (size >= 0),
  content_type TEXT        NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_documents_filename ON documents (filename);
```

## Running the Application

### Local Development

1.  **Install dependencies**:
    ```bash
    go mod tidy
    ```
2.  **Run the application**:
    ```bash
    go run cmd/api/main.go
    ```


### Using Docker (Single Container)

1.  **Build the image**:
    ```bash
    docker build -t godocapi .
    ```
2.  **Run the container**:
    ```bash
    docker run -p 8080:8080 --env-file .env godocapi
    ```

## Deployment

### Using Docker Compose

To start the full application stack (Frontend + Backend):

```bash
docker-compose up --build
```

The services will be available at:
-   **Frontend**: http://localhost:3000
-   **Backend**: http://localhost:8080

### Using Docker Swarm

1.  **Initialize Swarm** (if not already initialized):
    ```bash
    docker swarm init
    ```

2.  **Deploy the Stack**:
    ```bash
    docker stack deploy -c docker-stack.yml godocapi
    ```

3.  **Verify Services**:
    ```bash
    docker service ls
    ```

## API Documentation (Swagger)

The API documentation is auto-generated using Swagger. After running the application, you can access the Swagger UI at:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API Endpoints

### Health Check

-   `GET /api/v1/health`
    -   Check if the service is running.

### Documents

-   `POST /api/v1/documents`
    -   Upload a new document (Multipart Form: `file`).
-   `GET /api/v1/documents`
    -   List all documents.
-   `GET /api/v1/documents/:id`
    -   Get document metadata.
-   `GET /api/v1/documents/:id/download`
    -   Download the document file.
-   `DELETE /api/v1/documents/:id`
    -   Delete a document.
