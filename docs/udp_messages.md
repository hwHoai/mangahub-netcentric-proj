# UDP Server Messages (Quick Copy)

Dưới đây là danh sách các gói tin UDP JSON thường dùng để giao tiếp với UDP Notification Server.

### 1. Register Client (Device A - Receive Notifications)
Dùng để đăng ký địa chỉ UDP của client với server để nhận thông báo. Cần có `token` (JWT) để xác thực.
```json
{"action": "chapter:req_client_register", "token": "ACCESS_TOKEN_HERE", "payload": {"user_id": "65ebbe3b-df54-4d0f-b49a-e48ea6f98322"}}
```

### 2. Broadcast New Chapter Notification (From API/Service)
Sử dụng `HANDSHAKE_KEY` để xác thực.
```json
{"action": "chapter:impl_broadcast_notification", "token": "HANDSHAKE_KEY_HERE", "payload": {"manga_id": "29b2a759-ab30-4e71-a149-3e853756ce75", "chapter_id": "d2f5ea0b-9e2f-4d08-aaed-5166e047e1a7", "chapter_title": "Dog & Chainsaw", "chapter_number": 1}}
```

### 3. Broadcast New Chat Message Notification (From WebSocket Server)
Sử dụng `HANDSHAKE_KEY` để xác thực.
```json
{"action": "chat:impl_broadcast_message", "token": "HANDSHAKE_KEY_HERE", "payload": {"room_id": "29b2a759-ab30-4e71-a149-3e853756ce75", "sender_name": "Denji", "content": "Hello everyone!"}}
```

### 4. Sync Public Key (System Init)
```json
{"action": "pub_key:impl_sync_public_key", "token": "HANDSHAKE_KEY_HERE", "payload": {"public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"}}
```

### 5. Notification Received by Client (Example)
Khi có chương mới, server sẽ gửi về client đã đăng ký:
```json
{
  "action": "chapter:on_new_chapter_notification",
  "payload": {
    "manga_id": "29b2a759-ab30-4e71-a149-3e853756ce75",
    "chapter_id": "d2f5ea0b-9e2f-4d08-aaed-5166e047e1a7",
    "chapter_title": "Dog & Chainsaw",
    "chapter_number": 1
  }
}
```

---

## Example Test Workflow with `nc` (UDP mode)

1. **Terminal 1 (Listen for notifications)**:
   Mở một cổng UDP để nhận data (giả lập client):
   ```bash
   nc -u -l 9999
   ```

2. **Terminal 2 (Register Client)**:
   Gửi gói tin register tới UDP Server (mặc định port 8083):
   ```bash
   echo '{"action": "chapter:req_client_register", "token": "ACCESS_TOKEN_HERE", "payload": {"user_id": "USER_ID_HERE"}}' | nc -u -w1 localhost 8083
   ```

3. **Terminal 3 (Trigger Notification)**:
   Gửi gói tin broadcast (giả lập API Server):
   ```bash
   echo '{"action": "chapter:impl_broadcast_notification", "token": "HANDSHAKE_KEY_HERE", "payload": {"manga_id": "MANGA_ID_HERE", "chapter_id": "CHAPTER_ID_HERE", "chapter_title": "Test", "chapter_number": 1}}' | nc -u -w1 localhost 8083
   ```
