package web

import (
	"manajemen_tugas_master/model/domain"
)

type UpdateResponse struct {
	NameTask            string `json:"name_task,omitempty"`
	PlanningDescription string `json:"planning_description,omitempty"`
	PlanningStatus      string `json:"planning_status,omitempty"`
	ProjectStatus       string `json:"project_status,omitempty"`
	PlanningDueDate     string `json:"planning_due_date,omitempty"`
	ProjectDueDate      string `json:"project_due_date,omitempty"`
	Priority            string `json:"priority,omitempty"`
	ProjectComment      string `json:"project_comment,omitempty"`
	Manager             struct {
		ID     uint64 `json:"id,omitempty"`
		Email  string `json:"email,omitempty"`
		UserID uint64 `json:"user_id,omitempty"`
	} `json:"manager,omitempty"`
	Employee struct {
		ID     uint64 `json:"id,omitempty"`
		Email  string `json:"email,omitempty"`
		UserID uint64 `json:"user_id,omitempty"`
	} `json:"employee,omitempty"`
	PlanningFile struct {
		ID       uint64 `json:"id,omitempty"`
		FileUrl  string `json:"file_url,omitempty"`
		FileName string `json:"file_name,omitempty"`
	} `json:"planning_file,omitempty"`
	ProjectFile struct {
		ID       uint64 `json:"id,omitempty"`
		FileUrl  string `json:"file_url,omitempty"`
		FileName string `json:"file_name,omitempty"`
	} `json:"project_file,omitempty"`
}

func CreateResponseTask(taskModel *domain.Task) WebResponse {
	return WebResponse{
		Code:    200,
		Message: "Success",
		Data: domain.Task{

			ID:                  taskModel.ID,
			OwnerID:             taskModel.OwnerID,
			Owner:               taskModel.Owner,
			Manager:             taskModel.Manager,
			Employee:            taskModel.Employee,
			NameTask:            taskModel.NameTask,
			PlanningDescription: taskModel.PlanningDescription,
			PlanningFile:        taskModel.PlanningFile,
			PlanningStatus:      taskModel.PlanningStatus,
			ProjectFile:         taskModel.ProjectFile,
			ProjectStatus:       taskModel.ProjectStatus,
			PlanningDueDate:     taskModel.PlanningDueDate,
			ProjectDueDate:      taskModel.ProjectDueDate,
			Priority:            taskModel.Priority,
			ProjectComment:      taskModel.ProjectComment,
		},
	}
}

func CreateResponseTasks(tasksModel []*domain.Task) []WebResponse {
	var response []WebResponse
	for _, taskModel := range tasksModel {
		response = append(response, WebResponse{
			Code:    200,
			Message: "Success",
			Data: domain.Task{
				BoardID:             taskModel.BoardID,
				ID:                  taskModel.ID,
				OwnerID:             taskModel.OwnerID,
				Owner:               taskModel.Owner,
				Manager:             taskModel.Manager,
				Employee:            taskModel.Employee,
				NameTask:            taskModel.NameTask,
				PlanningDescription: taskModel.PlanningDescription,
				PlanningFile:        taskModel.PlanningFile,
				PlanningStatus:      taskModel.PlanningStatus,
				ProjectFile:         taskModel.ProjectFile,
				ProjectStatus:       taskModel.ProjectStatus,
				PlanningDueDate:     taskModel.PlanningDueDate,
				ProjectDueDate:      taskModel.ProjectDueDate,
				Priority:            taskModel.Priority,
				ProjectComment:      taskModel.ProjectComment,
			},
		})
	}
	return response
}
