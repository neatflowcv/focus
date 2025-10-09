package domain

type ExtraID string

type Extra struct {
	id       ExtraID
	parentID ExtraID
	leaf     bool
	status   TaskStatus
}

func NewExtra(
	id ExtraID,
	parentID ExtraID,
	leaf bool,
	status TaskStatus,
) *Extra {
	return &Extra{
		id:       id,
		parentID: parentID,
		leaf:     leaf,
		status:   status,
	}
}

func (e *Extra) ID() ExtraID {
	return e.id
}

func (e *Extra) Leaf() bool {
	return e.leaf
}

func (e *Extra) ParentID() ExtraID {
	return e.parentID
}

func (e *Extra) Status() TaskStatus {
	return e.status
}

func (e *Extra) IsCompleted() bool {
	return e.status == TaskStatusDone
}

func (e *Extra) SetLeaf(leaf bool) *Extra {
	ret := e.clone()
	ret.leaf = leaf
	ret.validate()

	return ret
}

func (e *Extra) SetParentID(parentID ExtraID) *Extra {
	ret := e.clone()
	ret.parentID = parentID
	ret.validate()

	return ret
}

func (e *Extra) SetStatus(status TaskStatus) *Extra {
	ret := e.clone()
	ret.status = status
	ret.validate()

	return ret
}

func (e *Extra) clone() *Extra {
	return NewExtra(e.id, e.parentID, e.leaf, e.status)
}

func (e *Extra) validate() {
	e.status.validate()
}
