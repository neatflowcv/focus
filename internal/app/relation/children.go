package relation

type ListChildrenInput struct {
	ParentID string
}

type ListChildrenOutput struct {
	IDs []string
}
