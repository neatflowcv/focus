package domain

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

func (s TaskStatus) validate() { //nolint:unused
	switch s {
	case TaskStatusTodo, TaskStatusDoing, TaskStatusDone:
		return
	default:
		panic("invalid task status: " + string(s))
	}
}
