package extra

type CreateExtraInput struct {
	ID       string
	ParentID string
}

type CreateExtraOutput struct {
	Extra Extra
}

type DeleteExtraInput struct {
	ID string
}

type ListExtrasInput struct {
	IDs []string
}

type ListExtrasOutput struct {
	Extras []*Extra
}

type Extra struct {
	Leaf   bool
	Status string
}
