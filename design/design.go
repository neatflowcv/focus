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

	dsl.Error("Unauthorized", dsl.ErrorResult, "Unauthorized")
	dsl.Error("InternalServerError", dsl.ErrorResult, "Internal server error")
	dsl.Error("TaskNotFound", dsl.ErrorResult, "Task not found")

	dsl.Method("setup", func() {
		dsl.Description("Setup the task service.")

		dsl.Payload(SetupTaskInput)

		dsl.HTTP(func() {
			dsl.POST("/setup")

			dsl.Header("authorization", dsl.String, "The authorization header")

			dsl.Response(dsl.StatusOK)
			dsl.Response("Unauthorized", dsl.StatusUnauthorized)
			dsl.Response("InternalServerError", dsl.StatusInternalServerError)
		})
	})

	dsl.Method("create", func() {
		dsl.Description("Create a new task.")

		dsl.Payload(CreateTaskInput)
		dsl.Result(CreateTaskOutput)

		dsl.HTTP(func() {
			dsl.POST("/")

			dsl.Header("authorization", dsl.String, "The authorization header")

			dsl.Response(dsl.StatusOK)
			dsl.Response("Unauthorized", dsl.StatusUnauthorized)
			dsl.Response("InternalServerError", dsl.StatusInternalServerError)
		})
	})

	dsl.Method("list", func() {
		dsl.Description("List all tasks.")

		dsl.Payload(func() {
			dsl.Attribute("authorization", dsl.String, "The authorization header")
			dsl.Attribute("parent_id", dsl.String, "The ID of the parent task")
			dsl.Attribute("recursive", dsl.Boolean, "Whether to include all subtasks recursively")

			dsl.Required("authorization")
		})
		dsl.Result(dsl.CollectionOf(TaskDetail))

		dsl.HTTP(func() {
			dsl.GET("/")

			dsl.Header("authorization", dsl.String, "The authorization header")
			dsl.Param("parent_id")
			dsl.Param("recursive")

			dsl.Response(dsl.StatusOK)
			dsl.Response("Unauthorized", dsl.StatusUnauthorized)
			dsl.Response("InternalServerError", dsl.StatusInternalServerError)
		})
	})

	dsl.Method("update", func() {
		dsl.Description("Update a task.")

		dsl.Payload(TaskUpdateInput)
		dsl.Result(TaskDetail)

		dsl.HTTP(func() {
			dsl.PATCH("/{task_id}")

			dsl.Header("authorization", dsl.String, "The authorization header")

			dsl.Response(dsl.StatusOK)
			dsl.Response("Unauthorized", dsl.StatusUnauthorized)
			dsl.Response("TaskNotFound", dsl.StatusNotFound)
			dsl.Response("InternalServerError", dsl.StatusInternalServerError)
		})
	})

	dsl.Method("delete", func() {
		dsl.Description("Delete a task.")

		dsl.Payload(TaskDeleteInput)

		dsl.HTTP(func() {
			dsl.DELETE("/{task_id}")

			dsl.Header("authorization", dsl.String, "The authorization header")

			dsl.Response(dsl.StatusNoContent)
			dsl.Response("Unauthorized", dsl.StatusUnauthorized)
			dsl.Response("TaskNotFound", dsl.StatusNotFound)
			dsl.Response("InternalServerError", dsl.StatusInternalServerError)
		})
	})
})

var CreateTaskInput = dsl.Type("CreateTaskInput", func() { //nolint:gochecknoglobals
	dsl.Attribute("authorization", dsl.String, "The authorization header")
	dsl.Attribute("parent_id", dsl.String, "The parent ID of the task")
	dsl.Attribute("title", dsl.String, "The title of the task")

	dsl.Required("authorization", "title")
})

var CreateTaskOutput = dsl.ResultType("CreateTaskOutput", func() { //nolint:gochecknoglobals
	dsl.Attribute("id", dsl.String, "The ID of the task")
	dsl.Attribute("created_at", dsl.Int64, "The timestamp when the task was created")

	dsl.Required("id", "created_at")
})

var TaskUpdateInput = dsl.Type("TaskUpdateInput", func() { //nolint:gochecknoglobals
	dsl.Attribute("authorization", dsl.String, "The authorization header")
	dsl.Attribute("task_id", dsl.String, "The ID of the task")

	dsl.Attribute("title", dsl.String, "The title of the task")
	dsl.Attribute("parent_id", dsl.String, "The parent ID of the task")
	dsl.Attribute("next_id", dsl.String, "The next ID of the task")
	dsl.Attribute("status", dsl.String, "The status of the task")
	dsl.Attribute("estimated_time", dsl.Int64, "The estimated time of the task")

	dsl.Required("authorization", "task_id")
})

var TaskDetail = dsl.ResultType("TaskDetail", func() { //nolint:gochecknoglobals
	dsl.Attribute("id", dsl.String, "The ID of the task")
	dsl.Attribute("parent_id", dsl.String, "The parent ID of the task")

	dsl.Attribute("title", dsl.String, "The title of the task")
	dsl.Attribute("status", dsl.String, "The status of the task")

	dsl.Attribute("is_leaf", dsl.Boolean, "Whether the task is a leaf task")

	// 시간 관련 속성
	dsl.Attribute("created_at", dsl.Int64, "The timestamp when the task was created")
	dsl.Attribute("completed_at", dsl.Int64, "The timestamp when the task was completed")
	dsl.Attribute("started_at", dsl.Int64, "The timestamp when the task was started")
	dsl.Attribute("lead_time", dsl.Int64, "The lead time of the task")
	dsl.Attribute("estimated_time", dsl.Int64, "The estimated time of the task")
	dsl.Attribute("actual_time", dsl.Int64, "The actual time of the task")

	dsl.Required("id", "title", "created_at", "status")
})

var TaskDeleteInput = dsl.Type("TaskDeleteInput", func() { //nolint:gochecknoglobals
	dsl.Attribute("authorization", dsl.String, "The authorization header")
	dsl.Attribute("task_id", dsl.String, "The ID of the task")

	dsl.Required("authorization", "task_id")
})

var SetupTaskInput = dsl.Type("SetupTaskInput", func() { //nolint:gochecknoglobals
	dsl.Attribute("authorization", dsl.String, "The authorization header")

	dsl.Required("authorization")
})
