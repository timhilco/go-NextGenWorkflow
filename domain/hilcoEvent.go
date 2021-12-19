package domain

// HilcoEvent is an type of hilcoEvent
type HilcoEvent struct {
	EventID string
}

// BakingEvent is an type of hilcoEvent
type BakingEvent struct {
	EventID string
	Task    string
}

// EventProcessor is an interface for processing events
type EventProcessor interface {
}
