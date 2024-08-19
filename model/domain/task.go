package domain

type Task struct {
	ID                  uint64         `json:"id" gorm:"primaryKey"`
	BoardID             uint64         `json:"board_id"`
	Board               Board          `json:"-" gorm:"foreignKey:BoardID;references:ID"`
	OwnerID             uint64         `json:"-"`
	Owner               Owner          `json:"owner" gorm:"foreignKey:OwnerID;references:ID"`
	Manager             []Manager      `json:"manager" gorm:"many2many:task_managers"`
	Employee            []Employee     `json:"employee" gorm:"many2many:task_employees"`
	NameTask            string         `json:"name_task" gorm:"size:255"`
	PlanningDescription string         `json:"planning_description"`
	PlanningFile        []PlanningFile `json:"planning_file"  gorm:"many2many:task_planning_files"`
	PlanningStatus      string         `json:"planning_status" gorm:"type:enum('Approved','Not Approved');default:'Not Approved'"`
	ProjectFile         []ProjectFile  `json:"project_file"  gorm:"many2many:task_project_files"`
	ProjectStatus       string         `json:"project_status" gorm:"type:enum('Done','Undone','Working');default:'Working'"`
	PlanningDueDate     string         `json:"planning_due_date" gorm:"size:255"`
	ProjectDueDate      string         `json:"project_due_date" gorm:"size:255"`
	Priority            string         `json:"priority" gorm:"type:enum('High','Medium','Low');default:'Medium'"`
	ProjectComment      string         `json:"project_comment"`
}
