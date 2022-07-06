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
