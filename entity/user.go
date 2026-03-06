package entity

type User struct {
	ID       int    `json:"id" gorm:"type:int;primary_key;auto_increment"`
	Username string `json:"username" gorm:"type:varchar(255);unique_index"`
	Name     string `json:"name" gorm:"type:varchar(255)"`
	Age      int    `json:"age" gorm:"type:int"`
}
