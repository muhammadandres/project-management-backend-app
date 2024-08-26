package web

import (
	"manajemen_tugas_master/model/domain"
)

type UpdateResponse struct {
	NameTask                  string `json:"name_task,omitempty"`
	PlanningDescriptionPersen string `json:"planning_description_persen,omitempty"`
	PlanningDescriptionFile   struct {
		ID       uint64 `json:"id,omitempty"`
		FileUrl  string `json:"file_url,omitempty"`
		FileName string `json:"file_name,omitempty"`
	} `json:"planning_description_file,omitempty"`
	PlanningStatus  string `json:"planning_status,omitempty"`
	ProjectStatus   string `json:"project_status,omitempty"`
	PlanningDueDate string `json:"planning_due_date,omitempty"`
	ProjectDueDate  string `json:"project_due_date,omitempty"`
	Priority        string `json:"priority,omitempty"`
	ProjectComment  string `json:"project_comment,omitempty"`
	Manager         struct {
		ID               uint64 `json:"id,omitempty"`
		Email            string `json:"email,omitempty"`
		UserID           uint64 `json:"user_id,omitempty"`
		InvitationStatus string `json:"invitation_status,omitempty"`
		InvitationID     uint64 `json:"invitation_id,omitempty"`
	} `json:"manager,omitempty"`
	Employee struct {
		ID               uint64 `json:"id,omitempty"`
		Email            string `json:"email,omitempty"`
		UserID           uint64 `json:"user_id,omitempty"`
		InvitationStatus string `json:"invitation_status,omitempty"`
		InvitationID     uint64 `json:"invitation_id,omitempty"`
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
	EmailsSent []string `json:"emails_sent,omitempty"`
}

func CreateResponseTask(taskModel *domain.TaskWithInvitation) WebResponse {
	return WebResponse{
		Code:    200,
		Message: "Success",
		Data:    taskModel,
	}
}

func CreateResponseTasks(tasksModel []*domain.TaskWithInvitation) []WebResponse {
	var response []WebResponse
	for _, taskModel := range tasksModel {
		response = append(response, WebResponse{
			Code:    200,
			Message: "Success",
			Data:    taskModel,
		})
	}
	return response
}
