package domain

import "time"

// Person represents a persons
type Person struct {
	InternalID string
	ExternalID string
	LastName   string
	FirstName  string
	BirthDate  time.Time
	Gender     string
}
