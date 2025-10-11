package domain

func SortTasks(tasks []*Task, parentID TaskID) []*Task {
	taskMap := make(map[TaskID]*Task)
	for _, task := range tasks {
		taskMap[task.ID()] = task
	}

	id := TaskDummyID(parentID)

	var ret []*Task

	task := taskMap[id]
	for task.NextID() != "" {
		next := taskMap[task.NextID()]
		ret = append(ret, next)
		task = next
	}

	return ret
}
