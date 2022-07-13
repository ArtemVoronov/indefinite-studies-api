package entities

type Tag struct {
	Id    int
	Name  string
	State string
}

const (
	TAG_STATE_NEW     string = "NEW"
	TAG_STATE_BLOCKED string = "BLOCKED"
	TAG_STATE_DELETED string = "DELETED"
)

func GetPossibleTagStates() []string {
	return []string{TAG_STATE_NEW, TAG_STATE_BLOCKED, TAG_STATE_DELETED}
}
