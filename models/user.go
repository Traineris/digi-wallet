package models

type User struct {
	GormModel
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Level    int    `json:"level" gorm:"default:1;not null"`
	Balance  int    `json:"balance" gorm:"default:0"`
}

func (User) TableName() string {
	return "users"
}
