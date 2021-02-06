package models

type UserIF interface {
	User() *User       // getter
	SetUser(*User)     // setter
	TableName() string // gormで使用するテーブルを強制指定
	Register() error
	GetByLoginID(loginID string) (*User, error)
}

type User struct {
	LoginID  string `gorm:"column:login_id"`
	Password string `gorm:"column:password"`
}
