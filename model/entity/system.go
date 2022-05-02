package entity

const SystemActionTable = "system_action"

type SystemAction struct {
	KeyName string `json:"key_name" gorm:"primaryKey"`
	Title   string `json:"title"`
	Sort    int    `json:"sort"`
}

func (*SystemAction) TableName() string { return SystemActionTable }
