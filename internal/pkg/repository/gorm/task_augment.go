package gorm

type TaskAugment struct {
	TaskID             string `gorm:"primaryKey"`
	Leaf               bool
	AllDescendantsDone bool
	TotalEstimatedTime int64
	TotalActualTime    int64
}
