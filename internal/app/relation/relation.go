package relation

type CreateRelationInput struct {
	ID       string
	ParentID string
}

type GetRelationInput struct {
	ID string
}

type GetRelationOutput struct {
	ID       string
	ParentID string
}
