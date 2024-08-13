package domain

type TaskWithInvitation struct {
	Task
	ManagersWithInvitation  []ManagerWithInvitation  `json:"manager"`
	EmployeesWithInvitation []EmployeeWithInvitation `json:"employee"`
}

type ManagerWithInvitation struct {
	Manager
	InvitationStatus string `json:"invitation_status,omitempty"`
}

type EmployeeWithInvitation struct {
	Employee
	InvitationStatus string `json:"invitation_status,omitempty"`
}
