package entities

import "time"

type Comment struct {
	Id               int
	Text             string
	UserId           int
	NoteId           int
	LinkdedCommentId int
	State            string
	CreateDate       time.Time
	LastUpdateDate   time.Time
}

const (
	COMMENT_STATE_NEW     string = "NEW"
	COMMENT_STATE_BLOCKED string = "BLOCKED"
	COMMENT_STATE_DELETED string = "DELETED"
)
