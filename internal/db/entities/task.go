package entities

type Task struct {
	Id    int
	Name  string
	State string
}

const (
	TASK_STATE_NEW     string = "NEW"
	TASK_STATE_DONE    string = "DONE"
	TASK_STATE_DELETED string = "DELETED"
)

func GetPossibleTaskStates() []string {
	return []string{TASK_STATE_NEW, TASK_STATE_DONE, TASK_STATE_DELETED}
}
