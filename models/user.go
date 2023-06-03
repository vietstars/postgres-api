package models

import (
  "gorm.io/gorm"
  "golang.org/x/crypto/bcrypt"
  "time"
)

type User struct {
  UserID uint `gorm:"primary_key;column:id" json:"user_id" csv:"id"`
  UserName string `gorm:"column:name" json:"name" csv:"name"`
  Email string `gorm:"typevarchar(100);unique_index;not null; column:email" json:"email" csv:"email"`
  Password string `json:"-" csv:"password"`
  Store []Store `gorm:"foreignKey:OwnerID;references:UserID" json:"store"`
  UserStatus uint8 `gorm:"column:status" json:"user_status" csv:"status"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by_id" csv:"created_by"`
  CreatedAt time.Time `gorm:"column:created_at" json:"created_at" csv:"created_at"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by_id" csv:"updated_by"`
  UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" csv:"updated_at"`
  DeletedAt gorm.DeletedAt `json:"-" csv:"deleted_at"`
  Version uint `gorm:"default:1" json:"version" csv:"version"`
}

type UserList []*User

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  if hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10);err == nil {
    tx.Statement.SetColumn("Password", string(hash))
  }
  tx.Statement.SetColumn("CreatedAt", time.Now())

  return 
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
  // if tx.Statement.Changed() { }
    tx.Statement.SetColumn("Version", u.Version+1)
    tx.Statement.SetColumn("UpdatedAt", time.Now())

  return 
}
