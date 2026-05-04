package enums

type ReadingStatus string

const (
	ReadingStatusNotStarted ReadingStatus = "not_started"
	ReadingStatusInProgress ReadingStatus = "in_progress"
	ReadingStatusCompleted  ReadingStatus = "completed"
)