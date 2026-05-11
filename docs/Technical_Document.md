# Technical Documentation — MangaHub

MangaHub is a high-performance, distributed content tracking platform designed for real-time progress synchronization and social interaction. This document provides an architectural deep-dive, setup instructions, and comprehensive API specifications for evaluators and developers.

---

## 1. Architecture Overview

MangaHub follows a **multi-protocol distributed architecture** centered around an **internal gRPC communication layer**. This design decoupling ensures that specific protocol spikes (e.g., a viral chat room on WebSocket) do not impact the core availability of the metadata API or synchronization layers.

### 1.1 Architecture Design
The MangaHub ecosystem is built as a distributed network of gateways communicating through an internal gRPC service layer.

![System Architecture](./docs/architecture.png)

### 1.2 System Component Diagram
The system is composed of five specialized services:

*   **`api-server` (HTTP Gateway)**: The primary entry point for web/mobile clients. Manages JWT-based identity and provides RESTful access to the manga library.
*   **`grpc-server` (Core Logic)**: It orchestrates all database persistence (SQLite via GORM), serving as internal service provider for all gateways.
*   **`tcp-server` (Sync Gateway)**: Manages persistent TCP connections for **cross-device reading synchronization**. When a user reads a chapter on one device, all other connected sessions are updated instantly.
*   **`udp-server` (Notification Gateway)**: A low-latency service for **real-time alerts** (new chapters, system announcements).
*   **`websocket-server` (Social Gateway)**: Facilitates room-based community chat using persistent WebSocket connections.

### 1.3 Database Design
The system uses a relational schema designed for efficient content tracking and user interaction data.

![Database Schema](./docs/db_diagram.png)

### 1.4 Core Design Principles & SOLID Architecture
MangaHub is designed with maintainability and modularity in mind, utilizing modern software engineering patterns:

*   **Clean Architecture & DDD**: Strict separation between delivery (`cmd/`) and business logic (`internal/`). Logic is partitioned into bounded contexts (`auth`, `manga`, `user`, etc.) to minimize coupling.
*   **Interface-Based Design (SOLID)**: 
    *   **Single Responsibility (SRP)**: Each component has a singular purpose—Gateways handle I/O, Services handle logic, and Repositories handle data.
    *   **Open/Closed (OCP)**: The **Dispatcher** system is open for extension (registering new action handlers) but closed for modification.
    *   **Liskov Substitution (LSP)**: All `impl/` packages satisfy their respective interfaces, allowing the system to swap implementations (e.g., different database providers) without affecting the high-level services.
    *   **Interface Segregation (ISP)**: Interfaces are lean and domain-specific (e.g., `MangaRepository` vs. `UserRepository`), ensuring clients only depend on methods they actually use.
    *   **Dependency Inversion (DIP)**: Every service and repository is defined by an **Interface**. High-level business logic depends on these abstractions rather than concrete types.
*   **Repository Pattern**: Encapsulates all GORM/SQLite interactions. This decouples the business layer from the persistence layer, allowing for database migrations (e.g., SQLite to PostgreSQL) without modifying core services.
*   **Specialized Socket Management**:
    *   **Dispatchers & Handlers**: For TCP/UDP/WS, incoming messages are routed via a **Dispatcher** to specific **Handlers** based on an `action` string. This prevents bloated "switch-case" blocks and keeps protocol logic modular.
    *   **Connection Pools**: Custom thread-safe managers handle client state, registration, and broadcasting. Stress tests demonstrated the ability to handle over 16,000 concurrent simulated sessions using Go channels and sync primitives.
*   **Security & Scale**: Cross-protocol authentication using RSA-signed JWTs, ensuring a secure experience across HTTP, TCP, and WebSocket.

---

## 2. Setup & Installation

### 2.1 Prerequisites
*   **Go**: 1.21, or higher.
*   **Protoc**: Protocol Buffer compiler (for generating gRPC code).
*   **Make**: Short CLI commands for task automation.

### 2.2 Installation Steps
1.  **Clone & Enter Repository**:
    ```bash
    git clone https://github.com/hwHoai/mangahub-netcentric-proj.git
    cd mangahub-netcentric-proj
    ```
2.  **Configuration**:
    ```bash
    cp .env.example .env
    ```
3.  **Install Dependencies**:
    ```bash
    go mod download
    ```
4.  **Optional: Seed Mock Data**:
    To populate the database with real manga data from MangaDex, uncomment line 47 in `cmd/grpc-server/main.go` before running the services:
    ```go
    // cmd/grpc-server/main.go:47
    seedData(dbConn) 
    ```
    *Note: The seeder only runs if the database is empty.*

### 2.3 Execution
**All-in-One Start (Recommended for Windows):**
```powershell
./run-all.ps1
```
**Or if Make are installed:**
```bash
make run-all
```

**Individual Services (via Makefile):**
```bash
make run-grpc  # Start core first
make run-api   # Start gateway
make run-tcp   # Start sync server
make run-udp   # Start notification server
make run-ws    # Start chat server
```

---

## 3. API & Protocol Specifications

### 3.1 HTTP (REST) API
**Port**: `8081` | **Base Path**: `/api/v1`

| Endpoint | Method | Auth | Description |
| :--- | :--- | :--- | :--- |
| `/health` | `GET` | No | System health check. |
| `/signup` | `POST` | No | Register a new account. |
| `/login` | `POST` | No | Login and obtain JWT tokens. |
| `/auth/refresh` | `POST` | No | Refresh expired access tokens. |
| `/auth/me` | `GET` | **Yes** | Get current authenticated user profile. |
| `/mangas` | `GET` | No | List mangas with pagination (`limit`, `offset`). |
| `/mangas/:id` | `GET` | No | Get detailed metadata for a specific manga. |
| `/mangas/:id/chapters`| `GET` | No | List all chapters for a specific manga. |
| `/mangas/:id/messages`| `GET` | No | Get community chat history for a manga room. |
| `/mangas/:id/chapters`| `POST` | **Yes** | Sync/Create new chapter from MangaDex (Triggers **UDP Notification**). |
| `/chapters/:id` | `GET` | No | Publicly read a chapter. |
| `/user/mangas/following`| `GET` | **Yes** | List all mangas followed by the user. |
| `/user/mangas/:id/follow`| `POST` | **Yes** | Follow a manga for updates. |
| `/user/mangas/:id/follow`| `DELETE` | **Yes** | Unfollow a manga. |
| `/user/history` | `GET` | **Yes** | Retrieve user's reading history. |
| `/user/chapters/:id` | `GET` | **Yes** | Read a chapter (Triggers **TCP Sync** across devices). |
| `/quotes` | `GET` | No | Educational scraper for random motivational quotes. |

### 3.2 Real-time Protocol Details

#### TCP Sync Protocol (Port `8082`)
Handles real-time device synchronization using JSON payloads.
*   **Action**: `chapter_sync:req_register_client` — Registers a device for sync.
*   **Action**: `chapter_sync:impl_broadcast_read` — Internal request trigger handler when user starts reading.
*   **Action**: `chapter_sync:on_new_read_progress` — Fire event sync new chapter progress to other devices.
*   **Action**: `pub_key:impl_sync_public_key` — System message to sync RSA public keys.
*   *Detailed Specs*: [docs/tcp_messages.md](./docs/tcp_messages.md)

#### UDP Notification Protocol (Port `8083`)
High-speed notification dispatch for new content.
*   **Action**: `chapter:req_client_register` — Subscribes a UDP client to alerts.
*   **Action**: `chapter:impl_broadcast_chapter` — Internal request to trigger handler to broadcast a new chapter notification.
*   **Action**: `chapter:on_new_chapter_notification` — Fire event new chapter notification to clients.
*   **Action**: `chat:impl_broadcast_message` — Internal request to trigger handler to broadcast a new chat message notification.
*   **Action**: `chat:on_new_message_notification` — Fire event new chat message notification to clients.
*   **Action**: `pub_key:impl_sync_public_key` — System message to sync RSA public keys.
*   *Detailed Specs*: [docs/udp_messages.md](./docs/udp_messages.md)

#### WebSocket Chat Protocol (Port `8085`)
Social interactions within manga-specific rooms.
*   **Endpoint**: `/ws` (Requires `manga_id` and `token` as query parameters).
*   **Key Sync**: `POST /impl/sync-public-key` — Internal endpoint for cross-protocol RSA key synchronization.
*   *Detailed Specs*: [docs/ws_messages.md](./docs/ws_messages.md)

### 3.3 Internal gRPC Communication Layer
Defined in `/proto`, the gRPC interfaces ensure type-safe, low-latency communication between services. The gRPC server acts as the central data orchestrator:

*   **Repository Pattern**: All database interactions (SQLite/GORM) are encapsulated within the `internal/repository/` layer.
*   **Service Layer Integration**: gRPC services are initialized with a database connection and maintain their own repository instances (e.g., `MangaRepository`, `UserRepository`).
*   **Encapsulation**: GRPC services calls repository methods for persistence, ensuring a clean separation between network protocols and data access.
*   **Core Services**:
    *   `UserService`: Manages user accounts and authentication state.
    *   `MangaService`: Handles content metadata and library queries.
    *   `MessageService`: Manages community chat persistence.
    *   `SessionService`: Tracks active protocol sessions and JWT lifecycles.
    *   `ChapterService` & `UserMangaService`: Manages content sync and user progress.

---

## 4. Performance & Validation

### 4.1 Benchmarking Architecture
The system includes custom-built stress testing tools designed to validate performance stability under high concurrent benchmark scenarios.

*   **Tools Used**:
    *   **Custom Go Benchmarkers**: High-performance scripts utilizing `goroutines`, `sync.WaitGroup`, and `atomic` counters for precise measurement.
    *   **Hey**: Used for HTTP/REST throughput testing.
    *   **Prometheus + Grafana**: Integrated into servers to export real-time metrics (Connections, Throughput, Memory).
*   **Location**: All benchmarking logic is located in the root `/benchmarks` directory, decoupled from the main application logic.
*   **Test Methodology**:
    *   **Phase 1: Registration**: Spawns thousands of concurrent connections to the gateway (TCP/UDP/WS).
    *   **Phase 2: Sync Burst**: Uses a synchronization barrier (Go channels) to trigger a simultaneous load burst across all active connections.

### 4.2 Automated Testing & Code Quality
MangaHub follows a strict testing regimen to ensure logic correctness across the Clean Architecture layers.

*   **Unit Tests**: Written using the standard Go `testing` package and **Testify** for assertions. These tests use **Mock Objects** to isolate business logic from database and network dependencies.
*   **Location**: `_test.go` files are located alongside the implementation in `internal/` packages (e.g., `internal/user/impl/user_service_test.go`).
*   **Execution**:
    ```bash
    make test-v           # Runs all tests with verbose output
    make test-no-cache    # Forces a fresh test run bypassing the Go cache
    ```
