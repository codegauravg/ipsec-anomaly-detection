# IPSec Anomaly Detection System

A full-stack security monitoring application that detects suspicious IPsec negotiation failures, tracks them in real time, and presents them in a dashboard for investigation.

## Overview

This project implements a distributed anomaly detection workflow:

- Firewall or log producers emit IPSec events
- Events are placed into AWS SQS
- A Go worker consumes the queue and performs threshold-based anomaly detection
- Redis is used for sliding-window counting per source IP
- PostgreSQL stores durable anomaly records
- A Node.js API exposes anomaly data and management endpoints
- A Vue dashboard visualizes active threats and allows operators to mark incidents as resolved

This architecture is designed to demonstrate queue-driven processing, real-time security analytics, and a clean separation between data collection, detection, persistence, and presentation.

## System Architecture

```mermaid
flowchart LR
    A[IPSec Logs / Firewall Events] --> B[AWS SQS]
    B --> C[Go Worker]
    C --> D[Redis: per-IP counters]
    C --> E[PostgreSQL: anomaly records]
    E --> F[Node API]
    F --> G[Vue Dashboard]
    G --> H[Operator resolves incident]
    H --> E
```

## Components

### 1. Worker
Location: `worker/`

Language: Go

Responsibilities:
- Read messages from SQS
- Parse event payloads
- Filter only `IKE_NEGOTIATION_FAILED` events
- Count failures per source IP using Redis
- Trigger anomaly when a threshold is breached
- Store detected anomaly in PostgreSQL
- Delete processed SQS messages

Key logic:
- Threshold = 50 failures
- Window = 60 seconds
- Threat type = `BRUTE_FORCE_IKE`
- Severity = 8

### 2. API
Location: `api/`

Language: Node.js + Express

Responsibilities:
- Expose REST endpoints for anomalies
- Query anomaly records from PostgreSQL
- Return dashboard statistics
- Allow marking an anomaly as resolved

### 3. Web Dashboard
Location: `web/`

Language: Vue 3 + TypeScript + Vite

Responsibilities:
- Fetch active anomalies from the API
- Render them in a table
- Allow refresh and resolve actions

### 4. Data Stores
- Redis: Sliding-window counts and fast in-memory state
- PostgreSQL: Persistent threat records and state history
- AWS SQS: Asynchronous queue for event ingestion and buffering

### 5. Infrastructure
- Terraform in `infra/` provisions AWS networking, SQS, ElastiCache Redis, Aurora PostgreSQL, and EKS
- Kubernetes manifests in `k8s/` define API and worker deployments
- Docker Compose in the root configures local Redis and Postgres for development

## Tech Stack

- Go 1.26
- Node.js / Express
- Vue 3
- Redis
- PostgreSQL
- AWS SQS
- Terraform
- Kubernetes
- Docker / Docker Compose

## Key Business Logic

An anomaly is detected when:

- the event action is `IKE_NEGOTIATION_FAILED`
- the source IP is tracked in Redis
- the count for that source IP exceeds or reaches `50` within the `60s` window
- the record is inserted into the database as a detected anomaly

## Database Schema

The main table is `ipsec_anomalies`:

- `id`: UUID primary key
- `timestamp`: event timestamp
- `source_ip`: attacker IP
- `threat_type`: threat category
- `severity`: integer severity score
- `raw_payload`: full event JSON
- `resolved`: whether the incident is handled

Schema source: `docs/initTable.sql`

## API Endpoints

Base path: `/api`

- `GET /api/anomalies` — list anomalies with optional filters
- `GET /api/anomalies/stats` — returns summary metrics
- `PATCH /api/anomalies/:id/resolve` — resolves an anomaly

## Local Setup

### 1. Start local infrastructure

```bash
docker-compose up -d
```

This starts:
- Redis on `localhost:6379`
- PostgreSQL on `localhost:5432`

### 2. Initialize database

Run the SQL from `docs/initTable.sql` against the local PostgreSQL instance.

### 3. Start the API

```bash
cd api
npm install
npm run dev
```

### 4. Start the Web UI

```bash
cd web
npm install
npm run dev
```

### 5. Start the Go worker

Set the required environment variables:

```bash
export SQS_QUEUE_URL="..."
export REDIS_URL="localhost:6379"
export DATABASE_URL="postgres://admin:password@localhost:5432/ipsec"
```

Then run:

```bash
cd worker
go run ./cmd/worker
```

## Production Notes

This project is structured as a cloud-native security monitoring solution using AWS services and Kubernetes deployments. The infrastructure layer is intentionally designed for scaling and isolation between components.

## Interview Talking Points

This project demonstrates:

- Distributed event-driven processing
- Real-time anomaly detection using Redis counters
- Durable storage and operational analytics in PostgreSQL
- Queue-based decoupling and async processing using SQS
- Clean separation of responsibilities across Go, Node, and Vue
- Cloud-native design with Terraform and Kubernetes
- Security operations workflow for threat investigation and response

## Project Structure

```text
.
├── api/
│   ├── src/
│   ├── Dockerfile
│   └── package.json
├── docs/
│   ├── initTable.sql
│   ├── replace.go
│   ├── replacements-bkp.json
│   └── replacements.json
├── infra/
│   ├── main.tf
│   ├── outputs.tf
│   └── variables.tf
├── k8s/
│   ├── api-deployment.yaml
│   ├── configmaps.yaml
│   ├── secrets.yaml
│   ├── web-deployment.yaml
│   └── worker-deployment.yaml
├── web/
│   ├── src/
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── worker/
│   ├── cmd/
│   ├── internal/
│   └── go.mod
├── docker-compose.yml
├── README.md
├── HLD.md
├── LLD.md
└── .gitignore
```

## Summary

This system is an end-to-end example of modern security monitoring architecture: event ingestion, stateful anomaly scoring, persistent threat recording, and a user-interface for response operations.
