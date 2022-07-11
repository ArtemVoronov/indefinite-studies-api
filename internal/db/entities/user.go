package entities

import "time"

type User struct {
	Id             int
	Login          string
	Email          string
	Password       string
	Role           string
	State          string
	CreateDate     time.Time
	LastUpdateDate time.Time
}

const (
	USER_STATE_NEW       string = "NEW"
	USER_STATE_CONFRIMED string = "CONFIRMED"
	USER_STATE_BLOCKED   string = "BLOCKED"
	USER_STATE_DELETED   string = "DELETED"
)

const (
	USER_ROLE_OWNER    string = "OWNER"
	USER_ROLE_RESIDENT string = "RESIDNET"
	USER_ROLE_GI       string = "GI"
)
