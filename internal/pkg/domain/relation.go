package domain

type RelationID string

func DummyRelationID(id RelationID) RelationID {
	return id + "-dummy"
}

type Relation struct {
	id       RelationID
	parentID RelationID
	nextID   RelationID
	version  uint64
}

func NewRelation(id RelationID, parentID RelationID, nextID RelationID, version uint64) *Relation {
	if id == "" {
		panic("id is required")
	}

	if version == 0 {
		panic("version is required")
	}

	return &Relation{
		id:       id,
		parentID: parentID,
		nextID:   nextID,
		version:  version,
	}
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

func (r *Relation) Version() uint64 {
	return r.version
}

func (r *Relation) SetNextID(nextID RelationID) *Relation {
	return NewRelation(r.id, r.parentID, nextID, r.version)
}

func (r *Relation) SetParentID(parentID RelationID) *Relation {
	return NewRelation(r.id, parentID, r.nextID, r.version)
}
