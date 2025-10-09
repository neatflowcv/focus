package extra

type SetDoneInput struct {
	ID string
}

type SetDoingInput struct {
	ID string
}

type SetTodoInput struct {
	ID string
}

type UpdateParentInput struct {
	ID       string
	ParentID string
}
