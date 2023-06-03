package repositories

import (
  "github.com/vietstars/postgres-api/dto"
  "github.com/vietstars/postgres-api/models"
)

//db.Unscoped().Where("id = 2").Find(&users)
//
// func (u *User) AfterFind(tx *gorm.DB) (err error) {
//   if u.Version == 0 {
//     u.Version = 1
//   }
//   return
// }

func NewUser(new dto.UserNew) (user *models.User, err error) {
  user = &models.User{
    UserName: new.UserName,
    Email: new.Email,
    Password: new.Password,
  }

  tx := models.DB.Begin()

  if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return nil, err
  }

  tx.Commit()

  return user, nil
}


func GetAllUsers() (users *models.UserList, err error) {
  if err = models.DB.Find(&users).Error; err != nil{

    return nil, err
  }

  return users, nil
}

func GetUserById(id uint) (user *models.User, err error) {
  if err = models.DB.First(&user, id).Error; err != nil{

    return user, nil
  }

  return nil, err
}

func GetUserByEmail(email string) (user *models.User, err error) {
  err = models.DB.First(&user, "email = ?", email).Error

  return user, err
}


func UpdateUserById(id uint, edit dto.LangEdit) (lang *models.Lang, err error) {

  tx := models.DB.Begin()

  if err := tx.Error; err != nil {

      return nil, err
  }

  if err := models.DB.Where("id = ? And version = ?", id, edit.Version).First(&lang).Error; err != nil{
    tx.Rollback()

    return nil, err
  }

  lang.Locale = edit.Locale
  lang.Group = edit.Group
  lang.Key = edit.Key
  lang.Val = edit.Val

  tx.Save(&lang)
  tx.Commit()

  return lang, nil
}
