package flow

type CreateTaskInput struct {
	Username string
	ParentID string
	Title    string
}

type ListTasksInput struct {
	Username  string
	ParentID  string
	Recursive bool
}
