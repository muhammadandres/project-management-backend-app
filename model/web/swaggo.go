package web

// user
type SignupRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type TokenResponse struct {
	Message string `json:"message" example:"Signup successfully"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" example:"user@example.com"`
	ResetCode   string `json:"reset_code" example:"12345"`
	NewPassword string `json:"new_password" example:"password123"`
}

type GetUserByIDResponse struct {
	Code    int        `json:"code" example:"200"`
	Message string     `json:"message" example:"Success"`
	Data    UserDetail `json:"data"`
}

type GetAllUsersResponse struct {
	Code    int          `json:"code" example:"200"`
	Message string       `json:"message" example:"Success"`
	Data    []UserDetail `json:"data"`
}

type UserDetail struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
}

type UpdateUser struct {
	Email string `json:"email" example:"user@example.com"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Error Message"`
}

type SuccessResponse struct {
	Error string `json:"message" example:"Success Message"`
}

//board
type BoardResponse struct {
	BoardID   uint64 `json:"board_id" example:"1"`
	NameBoard string `json:"name_board" example:"My Project Board"`
	CreatedBy struct {
		UserID    uint64 `json:"user_id" example:"1"`
		UserEmail string `json:"user_email" example:"user@example.com"`
	} `json:"board_created_by"`
	Tasks []TaskInfo `json:"tasks"`
}

// TaskInfo represents the task information in a board
type TaskInfo struct {
	ID                        uint64                    `json:"id" example:"1"`
	BoardID                   uint64                    `json:"board_id" example:"1"`
	Owner                     Owner                     `json:"owner"`
	Manager                   []ManagerWithInvitation   `json:"manager"`
	Employee                  []EmployeeWithInvitation  `json:"employee"`
	NameTask                  string                    `json:"name_task" example:"Implement feature X"`
	PlanningDescriptionPersen string                    `json:"planning_description_persen" example:"50"`
	PlanningDescriptionFile   []PlanningDescriptionFile `json:"planning_description_files"`
	PlanningFile              []PlanningFile            `json:"planning_file"`
	PlanningStatus            string                    `json:"planning_status" example:"Approved"`
	ProjectFile               []ProjectFile             `json:"project_file"`
	ProjectStatus             string                    `json:"project_status" example:"Undone"`
	PlanningDueDate           string                    `json:"planning_due_date" example:"2023-12-31"`
	ProjectDueDate            string                    `json:"project_due_date" example:"2024-01-15"`
	Priority                  string                    `json:"priority" example:"High"`
	ProjectComment            string                    `json:"project_comment" example:"Making good progress"`
}

type Owner struct {
	UserID uint64 `json:"user_id" example:"1"`
	Email  string `json:"email" example:"owner@example.com"`
}

type ManagerWithInvitation struct {
	UserID           uint64 `json:"user_id" example:"2"`
	Email            string `json:"email" example:"manager@example.com"`
	InvitationStatus string `json:"invitation_status" example:"Accepted"`
}

type EmployeeWithInvitation struct {
	UserID           uint64 `json:"user_id" example:"3"`
	Email            string `json:"email" example:"employee@example.com"`
	InvitationStatus string `json:"invitation_status" example:"Pending"`
}

type PlanningDescriptionFile struct {
	ID       uint64 `json:"id" example:"1"`
	FileName string `json:"file_name" example:"planning_description.pdf"`
	FileURL  string `json:"file_url" example:"https://bucket-name.s3.amazonaws.com/planning_description.pdf"`
}

type PlanningFile struct {
	ID       uint64 `json:"id" example:"2"`
	FileName string `json:"file_name" example:"planning_document.docx"`
	FileURL  string `json:"file_url" example:"https://bucket-name.s3.amazonaws.com/planning_document.docx"`
}

type ProjectFile struct {
	ID       uint64 `json:"id" example:"3"`
	FileName string `json:"file_name" example:"project_report.pdf"`
	FileURL  string `json:"file_url" example:"https://bucket-name.s3.amazonaws.com/project_report.pdf"`
}

type CreateBoardRequest struct {
	NameBoard string `json:"name_board" example:"New Project Board"`
}

type UpdateBoardRequest struct {
	NameBoard string `json:"name_board" example:"Updated Project Board"`
}

// task
type TaskCreateResponse struct {
	BoardID   uint64 `json:"board_id" example:"1"`
	TaskID    uint64 `json:"task_id" example:"1"`
	NameTask  string `json:"name_task" example:"Implement feature X"`
	OwnerID   uint64 `json:"owner_id" example:"1"`
	UserEmail string `json:"user_email" example:"user@example.com"`
	UserID    uint64 `json:"user_id" example:"1"`
}

type TaskResponse struct {
	ID                        uint64                    `json:"id" example:"1"`
	BoardID                   uint64                    `json:"board_id" example:"1"`
	Owner                     Owner                     `json:"owner"`
	Manager                   []ManagerWithInvitation   `json:"manager"`
	Employee                  []EmployeeWithInvitation  `json:"employee"`
	NameTask                  string                    `json:"name_task" example:"Implement feature X"`
	PlanningDescriptionPersen string                    `json:"planning_description_persen" example:"50"`
	PlanningDescriptionFile   []PlanningDescriptionFile `json:"planning_description_files"`
	PlanningFile              []PlanningFile            `json:"planning_file"`
	PlanningStatus            string                    `json:"planning_status" example:"Approved"`
	ProjectFile               []ProjectFile             `json:"project_file"`
	ProjectStatus             string                    `json:"project_status" example:"Undone"`
	PlanningDueDate           string                    `json:"planning_due_date" example:"2023-12-31"`
	ProjectDueDate            string                    `json:"project_due_date" example:"2024-01-15"`
	Priority                  string                    `json:"priority" example:"High"`
	ProjectComment            string                    `json:"project_comment" example:"Making good progress"`
}

type InvitationResponse struct {
	ID               uint64 `json:"id" example:"1"`
	TaskID           uint64 `json:"task_id" example:"1"`
	UserID           uint64 `json:"user_id" example:"2"`
	Role             string `json:"role" example:"Manager"`
	InvitationStatus string `json:"invitation_status" example:"Pending"`
}

type UpdateResponseTask struct {
	TaskID                    uint64             `json:"task_id" example:"1"`
	BoardID                   uint64             `json:"board_id" example:"1"`
	NameTask                  string             `json:"name_task" example:"Updated Task Name"`
	PlanningDescriptionPersen string             `json:"planning_description_persen" example:"75"`
	PlanningStatus            string             `json:"planning_status" enums:"Approved,Not Approved" example:"Approved"`
	ProjectStatus             string             `json:"project_status" enums:"Working,Done,Undone" example:"Working"`
	PlanningDueDate           string             `json:"planning_due_date" example:"2023-12-31"`
	ProjectDueDate            string             `json:"project_due_date" example:"2024-01-15"`
	Priority                  string             `json:"priority" enums:"Low,Medium,High" example:"High"`
	ProjectComment            string             `json:"project_comment" example:"Updated project comment"`
	Managers                  []ManagerResponse  `json:"managers"`
	Employees                 []EmployeeResponse `json:"employees"`
	PlanningFiles             []FileResponse     `json:"planning_files"`
	ProjectFiles              []FileResponse     `json:"project_files"`
	PlanningDescriptionFiles  []FileResponse     `json:"planning_description_files"`
	EmailsSent                string             `json:"emails_sent" example:"Email send infomation"`
}

type ManagerResponse struct {
	Email string `json:"email" example:"manager@example.com"`
}

type EmployeeResponse struct {
	Email string `json:"email" example:"employee@example.com"`
}

type FileResponse struct {
	FileUrl  string `json:"file_url" example:"https://example.com/file.pdf"`
	FileName string `json:"file_name" example:"document.pdf"`
}
