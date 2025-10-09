package domain

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

func validateTaskStatus(status TaskStatus) {
	switch status {
	case TaskStatusTodo, TaskStatusDoing, TaskStatusDone:
		return
	}

	panic("invalid task status: " + string(status))
}
