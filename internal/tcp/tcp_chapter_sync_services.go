package tcp_services

type TCPChapterSyncServices interface {
	SyncReading(userID string, chapterID string) error
}