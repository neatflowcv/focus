package domain

type RelationID string

func DummyRelationID(id RelationID) RelationID {
	return id + "-dummy"
}

type Relation struct {
	id       RelationID
	parentID RelationID
	nextID   RelationID
}

func NewRelation(id RelationID, parentID RelationID, nextID RelationID) *Relation {
	if id == "" {
		panic("id is required")
	}

	return &Relation{id: id, parentID: parentID, nextID: nextID}
}

func (r *Relation) ID() RelationID {
	return r.id
}

func (r *Relation) ParentID() RelationID {
	return r.parentID
}

func (r *Relation) NextID() RelationID {
	return r.nextID
}

func (r *Relation) SetNextID(nextID RelationID) *Relation {
	return NewRelation(r.id, r.parentID, nextID)
}
