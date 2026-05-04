# Features & Requirements - Net-centric Manga Hub

## 1. Authentication
- User can create an account with username and password.
- User can login with username and password.
- JWT-based authentication.

## 2. manga endpoint [search, get detail]
- Search manga by title.
- Get manga details including title, description, author, genres, chapters, etc.

## 3. user endpoint [update reading lib]

4. Viết TCP server lắng nghe nhiều kết nối đồng thời (dùng Goroutines). Nhận và Broadcast (phát sóng) tiến độ đọc truyện của user qua định dạng JSON. (chưa hiểu, đọc truyện thì mạnh ai nấy kéo truyện về đọc, mạnh ai nấy cập nhật tiến độ của mình, tại sao cần tcp làm gì)

5. Thông báo chapter mới [UDP notification, pub/sub]

6. chat [websocket]

7. Định nghĩa file .proto cho 2-3 services nội bộ. Code gRPC server thực hiện các hàm: GetManga, SearchManga, UpdateProgress. (chưa hiểu làm gì, microservices hả)

8. tạo wishlist riêng

9. Gợi ý truyện

10. Advanced filtering

11. Rating

# References
- [Database Schema](https://dbdiagram.io/d/Net-centric-mangahub-69e0f29d8089629684b4cd7c)