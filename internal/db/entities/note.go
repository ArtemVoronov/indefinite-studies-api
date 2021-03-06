package entities

import "time"

type Note struct {
	Id             int
	Text           string
	Topic          string
	TagId          int
	UserId         int
	State          string
	CreateDate     time.Time
	LastUpdateDate time.Time
}

const (
	NOTE_STATE_NEW     string = "NEW"
	NOTE_STATE_BLOCKED string = "BLOCKED"
	NOTE_STATE_DELETED string = "DELETED"
)

func GetPossibleNoteStates() []string {
	return []string{NOTE_STATE_NEW, NOTE_STATE_BLOCKED, NOTE_STATE_DELETED}
}
