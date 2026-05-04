# TCP Server Messages (Quick Copy)

Dưới đây là danh sách các gói tin TCP JSON thường dùng, được định dạng sẵn để copy.

### 1. Register Client (Device A)
Sử dụng JWT Token nhận được từ API Login/Signup.
```json
{"action": "chapter_sync:req_register_client", "payload": {"user_id": "65ebbe3b-df54-4d0f-b49a-e48ea6f98322"}, "token": "ACCESS_TOKEN_HERE"}
```

### 2. Broadcast Read Progress (From Server/API)
Sử dụng `HANDSHAKE_KEY` cấu hình trong `.env`.
```json
{"action": "chapter_sync:impl_broadcast_read", "payload": {"user_id": "65ebbe3b-df54-4d0f-b49a-e48ea6f98322", "chapter_id": "d2f5ea0b-9e2f-4d08-aaed-5166e047e1a7"}, "token": "HANDSHAKE_KEY_HERE"}
```

### 3. Sync Public Key (System Init)
```json
{"action": "pub_key:impl_sync_public_key", "payload": {"public_key": "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----\n"}, "token": "HANDSHAKE_KEY_HERE"}
```

### 4. Response From Server (What Client Receives)
Khi có thiết bị khác cập nhật tiến trình đọc, server sẽ gửi về:
```json
{"action":"chapter_sync:on_new_read_progress","payload":{"chapter_id":"d2f5ea0b-9e2f-4d08-aaed-5166e047e1a7"}}
```

---

## Example Test Workflow with `nc`

1. **Terminal 1 (Listen for changes)**:
   ```bash
   nc localhost 8082
   # Paste Register Client JSON above
   ```

2. **Terminal 2 (Trigger change)**:
   ```bash
   echo '{"action": "chapter_sync:impl_broadcast_read", "payload": {"user_id": "65ebbe3b-df54-4d0f-b49a-e48ea6f98322", "chapter_id": "d2f5ea0b-9e2f-4d08-aaed-5166e047e1a7"}, "token": "HANDSHAKE_KEY_HERE"}' | nc localhost 8082
   ```

3. **Check Terminal 1**: Bạn sẽ thấy message `on_new_read_progress` được đẩy về.
