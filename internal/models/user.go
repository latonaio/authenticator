package models

type UserIF interface {
	User() *User       // getter
	SetUser(*User)     // setter
	TableName() string // gormで使用するテーブルを強制指定
	Register() error
	GetByLoginID(loginID string) (*User, error)
}

type User struct {
	ID       int `gorm:"primaryKey"`
	LoginID  string
	Password string
	//CreatedAt time.Time  `gorm:"create_at"`
	//UpdatedAt time.Time  `gorm:"update_at"`
	//DeletedAt *time.Time `gorm:"delete_at"`
}
