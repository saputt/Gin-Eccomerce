package user

type UserRole string

const (
	RoleBuyer  UserRole = "BUYER"
	RoleSeller UserRole = "SELLER"
)

type User struct {
	Id       string   `json:"id" gorm:"primaryKey"`
	Name     string   `json:"name" gorm:"not null"`
	Email    string   `json:"email" gorm:"uniqueIndex;not null"`
	Password string   `json:"password"`
	Role     UserRole `json:"role" gorm:"not null;default:BUYER"`
}
