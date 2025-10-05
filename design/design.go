package design

import (
	"goa.design/goa/v3/dsl"
)

var _ = dsl.API("focus", func() {
	dsl.HTTP(func() {
		dsl.Path("/focus")
	})
})

var _ = dsl.Service("task", func() {
	dsl.HTTP(func() {
		dsl.Path("/tasks")
	})

	dsl.Method("create", func() {
		dsl.Description("Create a new task.")

		dsl.Payload(TaskInput)
		dsl.Result(TaskDetail)

		dsl.HTTP(func() {
			dsl.POST("/")

			dsl.Response(dsl.StatusOK)
		})
	})

	dsl.Method("list", func() {
		dsl.Description("List all tasks.")

		dsl.Payload(func() {
			dsl.Attribute("parent_id", dsl.String, "The ID of the parent task")
			dsl.Attribute("recursive", dsl.Boolean, "Whether to include all subtasks recursively")
		})
		dsl.Result(dsl.CollectionOf(TaskDetail))

		dsl.HTTP(func() {
			dsl.GET("/")

			dsl.Param("parent_id")
			dsl.Param("recursive")

			dsl.Response(dsl.StatusOK)
		})
	})
})

var TaskInput = dsl.Type("TaskInput", func() { //nolint:gochecknoglobals
	dsl.Attribute("parent_id", dsl.String, "The parent ID of the task")
	dsl.Attribute("title", dsl.String, "The title of the task")

	dsl.Required("title")
})

var TaskDetail = dsl.ResultType("TaskDetail", func() { //nolint:gochecknoglobals
	dsl.Attribute("id", dsl.String, "The ID of the task")
	dsl.Attribute("parent_id", dsl.String, "The parent ID of the task")
	dsl.Attribute("title", dsl.String, "The title of the task")
	dsl.Attribute("created_at", dsl.Int64, "The timestamp when the task was created")
	dsl.Attribute("status", dsl.String, "The status of the task")
	dsl.Attribute("order", dsl.Float64, "The order of the task")

	dsl.Required("id", "title", "created_at", "status", "order")
})
