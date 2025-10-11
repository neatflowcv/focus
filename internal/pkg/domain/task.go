package domain

import (
	"time"
)

type TaskID string

type Task struct {
	id       TaskID
	parentID TaskID
	nextID   TaskID

	title     string
	createdAt time.Time
	version   uint64
}

func NewTask(
	id TaskID,
	parentID TaskID,
	nextID TaskID,
	title string,
	createdAt time.Time,
	version uint64,
) *Task {
	ret := &Task{
		id:        id,
		parentID:  parentID,
		nextID:    nextID,
		title:     title,
		createdAt: createdAt,
		version:   version,
	}
	ret.validate()

	return ret
}

func (t *Task) Dummy() *Task {
	return &Task{
		id:        t.id + "-dummy",
		parentID:  t.id,
		nextID:    "",
		title:     "",
		createdAt: t.createdAt,
		version:   t.version,
	}
}

func (t *Task) Equals(other *Task) bool {
	return t.id == other.id &&
		t.parentID == other.parentID &&
		t.nextID == other.nextID &&
		t.title == other.title &&
		t.createdAt.Equal(other.createdAt) &&
		t.version == other.version
}

func (t *Task) ID() TaskID {
	return t.id
}

func (t *Task) ParentID() TaskID {
	return t.parentID
}

func (t *Task) NextID() TaskID {
	return t.nextID
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) Version() uint64 {
	return t.version
}

func (t *Task) SetParentID(parentID TaskID) *Task {
	ret := t.clone()
	ret.parentID = parentID
	ret.validate()

	return ret
}

func (t *Task) SetNextID(nextID TaskID) *Task {
	ret := t.clone()
	ret.nextID = nextID
	ret.validate()

	return ret
}

func (t *Task) SetTitle(title string) *Task {
	ret := t.clone()
	ret.title = title
	ret.validate()

	return ret
}

func (t *Task) clone() *Task {
	return NewTask(t.id, t.parentID, t.nextID, t.title, t.createdAt, t.version)
}

func (t *Task) validate() {
	if t.id == "" {
		panic("id is required")
	}

	if t.title == "" {
		panic("title is required")
	}

	if t.version == 0 {
		panic("version is required")
	}
}
