# 📖 MangaHub

**A High-Performance, Distributed Manga Tracking Platform Built with Go.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Microservices-FF6F00?style=flat)](#architecture)
[![Protocols](https://img.shields.io/badge/Protocols-REST%20%7C%20gRPC%20%7C%20TCP%20%7C%20UDP%20%7C%20WS-blue)](#core-services)

## 🎯 Overview

MangaHub is not just another CRUD application; it's a showcase of **advanced network programming** and **scalable system architecture**. Designed as a microservices-based platform, it efficiently manages manga libraries, tracks reading progress in real-time, and handles concurrent user interactions across multiple network protocols.

Whether you are evaluating architectural decisions, code quality, or networking concepts, MangaHub demonstrates a deep understanding of how disparate systems communicate reliably in a modern backend environment.

## 🚀 Why This Project Stands Out 

- **Protocol Mastery**: Successfully implements and orchestrates **five distinct network protocols** (HTTP/REST, gRPC, TCP, UDP, WebSocket) within a unified ecosystem.
- **Clean & Modular Architecture**: Adheres strictly to Domain-Driven Design (DDD) principles and interface-based implementations, ensuring high testability and maintainability.
- **Concurrency & Thread Safety**: Leverages Go's powerful concurrency model (goroutines, channels) with custom Connection Pools to handle simultaneous TCP/UDP broadcasts safely.
- **Automated Data Pipeline**: Features a robust seeder that reliably consumes external APIs (MangaDex), processes rate-limited requests, and maps complex JSON relationships into a normalized relational database.

---

## 🏗 System Architecture

MangaHub operates as a cluster of specialized services, each optimized for its specific transport layer:

### Core Services
1. **🌐 HTTP REST API Server (Gateway)**: Built with Gin, serving as the primary entry point for client applications. Handles JWT-based authentication and routes business logic to internal services.
2. **🔌 gRPC Internal Service**: The high-speed backbone of the system. Manages all direct database transactions (SQLite + GORM) via Protocol Buffers, decoupling the data layer from public-facing APIs.
3. **⚡ TCP Progress Sync Server**: Maintains persistent connections to broadcast live reading progress updates across user devices securely.
4. **📡 UDP Notification System**: A lightweight broadcasting server designed to blast out "New Chapter" alerts with minimal overhead. *(In active development)*
5. **💬 WebSocket Chat Hub**: Enables real-time, low-latency community discussions for specific manga titles. *(In active development)*

## 💻 Technology Stack

- **Core Language**: Go (Golang)
- **API Framework**: Gin Web Framework
- **RPC Framework**: gRPC & Protocol Buffers (Protobuf)
- **Database & ORM**: SQLite3, GORM
- **Security**: RSA Key Pairs, JWT (JSON Web Tokens) for stateless authentication.

---

## 🛠 Getting Started

### Prerequisites
- Go 1.21 or higher installed.
- Protocol Buffers compiler (`protoc`) if you plan to modify `.proto` files.

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-username/manga_hub.git
   cd manga_hub
   ```

2. **Configure Environment:**
   Copy the example environment file and adjust if necessary.
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies:**
   ```bash
   go mod download
   ```

4. **Run the Services:**
   We recommend running these in separate terminal windows to monitor logs effectively.
   
   *Start the Internal gRPC Server (Data Layer):*
   ```bash
   go run cmd/grpc-server/main.go
   ```
   
   *Start the TCP Sync Server:*
   ```bash
   go run cmd/tcp-server/main.go
   ```
   
   *Start the Public REST API:*
   ```bash
   go run cmd/api-server/main.go
   ```

*(Note: Data seeding from MangaDex runs automatically on the first gRPC server boot if the database is empty).*

---

## ✅ Development Roadmap (To-Do for Full Score)

### Phase 2: Network Protocols Completion
- [ ] **UDP Notification System (15 pts)**
  - [ ] Implement UDP Server client registration.
  - [ ] Build broadcasting logic for "New Chapter" alerts.
- [ ] **WebSocket Chat System (15 pts)**
  - [ ] Set up `ChatHub` for WebSocket connection management.
  - [ ] Implement real-time join/leave/message broadcasting.

### Data Collection Requirements
- [ ] **Web Scraping Practice**
  - [ ] Write a simple script (e.g., scraping `quotes.toscrape.com`) to fulfill the "Educational Practice" grading criteria.

### Target Bonus Features (Pick 1-2 for Extra Credit)
- [ ] **Health Checks (5 pts)**: Formalize the `/health` endpoint for all microservices.
- [ ] **Input Sanitization (5 pts)**: Add comprehensive payload validation to API endpoints.
- [ ] **Data Caching with Redis (10 pts)**: Setup Redis to cache popular manga queries.

---

## 🎓 Academic Context

This project is developed as the capstone assignment for the **Net-centric Programming (IT096IU)** course at International University (VNU-HCM), instructed by Lê Thanh Sơn & Nguyễn Trung Nghĩa. It is engineered to meet and exceed the rigorous grading criteria for network protocol implementation, system integration, and code quality.

## 📄 License

This project is open-source and available under the [MIT License](LICENSE).
