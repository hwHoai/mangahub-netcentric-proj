# MangaHub Run All Services Script (Windows)

Write-Host "🚀 Starting MangaHub Microservices..." -ForegroundColor Cyan

# 1. Start gRPC Server
Write-Host "-> Starting gRPC Server..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "make run-grpc"

# Wait a bit for gRPC to be ready
Start-Sleep -Seconds 2

# 2. Start TCP Server
Write-Host "-> Starting TCP Server..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "make run-tcp"

# 3. Start UDP Server
Write-Host "-> Starting UDP Server..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "make run-udp"

# 4. Start WebSocket Server
Write-Host "-> Starting WebSocket Server..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "make run-ws"

# 5. Start API Server (Gateway)
Write-Host "-> Starting API Server..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "make run-api"

Write-Host "✅ All services are launching in separate windows." -ForegroundColor Green
Write-Host "Please check each window for logs."
