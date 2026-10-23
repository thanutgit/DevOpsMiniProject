package entity

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"type:varchar(255);uniqueIndex"`
	Name     string `json:"name" gorm:"type:varchar(255)"`
	Age      int    `json:"age" gorm:"type:int"`
}

func (User) TableName() string {
	return "users"
}
