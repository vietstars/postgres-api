package models

import (
  "gorm.io/gorm"
  "golang.org/x/crypto/bcrypt"
  "time"
)

type User struct {
  ID uint `gorm:"primary_key" json:"id"`
  Email string `json:"email"`
  Password string `json:"-"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by_id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by_id"`
  UpdatedAt time.Time `json:"updated_at"` //gorm:"default:current_timestamp"
  DeletedAt gorm.DeletedAt `json:"-"`
  Version uint `json:"version"`
}

type Auth struct {
  Info *User `json:"info"`
  Token string `json:"token"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
  if hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10);err == nil {
    tx.Statement.SetColumn("Password", string(hash))
  }

  tx.Statement.SetColumn("Version", u.Version+1)
  tx.Statement.SetColumn("CreatedAt", time.Now())

  return 
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
  if tx.Statement.Changed() {
    tx.Statement.SetColumn("Version", u.Version+1)
    tx.Statement.SetColumn("UpdatedAt", time.Now())
  }

  return 
}

//db.Unscoped().Where("id = 2").Find(&users)
//
// func (u *User) AfterFind(tx *gorm.DB) (err error) {
//   if u.Version == 0 {
//     u.Version = 1
//   }
//   return
// }

func NewUser(email string, password string) (user *User, err error) {
  user = &User{
    Email: email,
    Password: password,
  }

  DB.Create(&user)

  return user, nil
}

func GetFirstById(id uint) (user *User, err error) {

  if err = DB.First(&user, id).Error; err != nil{

    return user, nil
  }

  return nil, err
}

func GetFirstByEmail(email string) (user *User, err error) {
  err = DB.First(&user, "email = ?", email).Error

  return user, err
}
