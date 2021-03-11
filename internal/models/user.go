package models

import "time"

type UserIF interface {
	User() *User       // getter
	SetUser(*User)     // setter
	TableName() string // gormで使用するテーブルを強制指定
	Register() error
	Login() error
	GetByLoginID(loginID string) (*User, error)
}

type User struct {
	ID          int       `gorm:"primaryKey"`
	LoginID     string    `gorm:"column:login_id"`
	Password    string    `gorm:"column:password"`
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time // column name is `created_at`
	UpdatedAt   time.Time // column name is `updated_at`
}
