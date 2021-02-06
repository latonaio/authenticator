package models

import "bitbucket.org/latonaio/authenticator/pkg/db"

func NewUser() UserIF {
	return &User{}
}

func (u *User) Register() error {
	result := db.ConPool.Con.Create(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) GetByLoginID(loginID string) (*User, error) {
	result := db.ConPool.Con.Model(u).Where("login_id = ?", loginID).First(u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (u *User) TableName() string {
	return db.ConPool.Info.TableName
}

func (u *User) User() *User {
	return u
}

func (u *User) SetUser(user *User) {
	u.LoginID = user.LoginID
	u.Password = user.Password
	//u.CreatedAt = user.CreatedAt
	//u.UpdatedAt = user.UpdatedAt
	//u.DeletedAt = user.DeletedAt
}
