package handlers

import (
	"net/http"
	api "test_task/gen/task"
	"test_task/service"

	"github.com/labstack/echo/v4"
)

type TaskAPI struct {
	manager *service.TaskManager
}

func NewTaskAPI(m *service.TaskManager) api.ServerInterface {
	return &TaskAPI{manager: m}
}

func (h *TaskAPI) CreateTask(c echo.Context) error {
	var body struct {
		payload interface{}
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := h.manager.Submit(c.Request().Context(), body.payload)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, err)
	}
	status := api.Pending
	return c.JSON(http.StatusAccepted, api.Task{
		Id:     &id,
		Status: &status,
	})
}

func (h *TaskAPI) GetTask(c echo.Context, id string) error {
	if t, ok := h.manager.GetID(id); ok {
		status := api.TaskStatus(t.Status)

		var result *map[string]interface{}
		switch v := t.Result.(type) {
		case map[string]interface{}:
			result = &v
		case map[string]string:
			converted := make(map[string]interface{}, len(v))
			for k, val := range v {
				converted[k] = val
			}
			result = &converted
		default:
			result = nil
		}

		return c.JSON(http.StatusOK, api.Task{
			Id:     &t.ID,
			Status: &status,
			Result: result,
			Error:  &t.Error,
		})
	}
	return c.JSON(http.StatusNotFound, api.Error{Message: "not found"})
}
