package extra

import "time"

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

type UpdateEstimatedTimeInput struct {
	ID            string
	EstimatedTime time.Duration
}

type ListExtrasInput struct {
	IDs []string
}

type ListExtrasOutput struct {
	Extras []Extra
}

type Extra struct {
	EstimatedTime time.Duration
	ActualTime    time.Duration
	StartedAt     time.Time
	Leaf          bool
}

type UpdateActualTimeInput struct {
	ID         string
	ActualTime time.Duration
}
