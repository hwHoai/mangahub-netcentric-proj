# WebSocket Server Messages (Quick Copy)

Dưới đây là danh sách các gói tin WebSocket JSON thường dùng để giao tiếp với WebSocket Chat Server.

### 1. Connection URL
Để kết nối vào một phòng chat (manga room), bạn cần truyền `manga_id` và `token` (JWT) qua query string.
```text
ws://localhost:8085/ws?manga_id=29b2a759-ab30-4e71-a149-3e853756ce75&token=ACCESS_TOKEN_HERE
```

### 2. Send Message (Client -> Server)
Sau khi kết nối thành công, client gửi tin nhắn bằng format:
```json
{
  "content": "Chào mọi người, mình mới tham gia phòng này!"
}
```

### 3. Receive Message (Server -> Client)
Khi có tin nhắn mới trong phòng, server sẽ broadcast tới tất cả client:
```json
{
  "room_id": "29b2a759-ab30-4e71-a149-3e853756ce75",
  "content": "Chào mọi người, mình mới tham gia phòng này!",
  "sender": "65ebbe3b-df54-4d0f-b49a-e48ea6f98322"
}
```

---

## Example Test Workflow with `wscat`

1. **Install wscat** (nếu chưa có):
   ```bash
   npm install -g wscat
   ```

2. **Connect to Server**:
   ```bash
   wscat -c "ws://localhost:8085/ws?manga_id=MANGA_ID_HERE&token=TOKEN_HERE"
   ```

3. **Send Message**:
   Trong terminal của `wscat`, paste:
   ```json
   {"content": "Hello world!"}
   ```

4. **Observe Response**:
   Bạn sẽ nhận được chính tin nhắn đó (broadcast) và các tin nhắn từ user khác trong cùng `manga_id`.
