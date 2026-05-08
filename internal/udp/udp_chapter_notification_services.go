package udp_services

// UDPChapterNotificationServices is the interface for sending UDP notifications.
type UDPChapterNotificationServices interface {
	SendNewChapterNotification(mangaID, chapterID, title string, chapterNumber float64) error
	SendNewMessageNotification(roomID, senderName, content string) error
	Close() error
}