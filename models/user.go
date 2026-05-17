package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Password string
	Name     string

}
func (User) TableName() string {
return "user"}