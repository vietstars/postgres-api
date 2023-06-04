package repositories

import (
  "github.com/vietstars/postgres-api/dto"
  "github.com/vietstars/postgres-api/models"
)

func NewLang(new dto.LangNew) (lang *models.LangEntity, err error) {
  lang = &models.LangEntity{
    Locale: new.Locale,
    Group: new.Group,
    Key: new.Key,
    Val: new.Val,
  }

  tx := models.DB.Table("langs").Begin()

  if err := tx.Create(&lang).Error; err != nil {
    tx.Rollback()
    return nil, err
  }

  tx.Commit()

  return lang, nil
}

func DelLangById(id uint, del dto.LangDel) (result bool, err error) {
  var lang models.LangEntity

  tx := models.DB.Table("langs").Begin()

  if err := tx.Error; err != nil {

      return  false, err
  }

  if err := models.DB.Where("id = ? And version = ?", id, del.Version).First(&lang).Error; err != nil{
    tx.Rollback()

    return false, err
  }

  if del.ForceDel {
    tx.Unscoped().Delete(&lang)
  } else {
    tx.Delete(&lang)
  }

  tx.Commit()

  return true, nil
}

func GetAllLang() (langs *models.LangListEntity, err error) {
  if err = models.DB.Find(&langs).Error; err != nil{

    return nil, err
  }

  return langs, nil
}

func GetLangsByLocale(lg string) (langs *models.LangListEntity, err error) {
  err = models.DB.Find(&langs, "locale = '_' OR locale = ?", lg).Error

  return langs, err
}

func GetLangById(id uint) (lang *models.LangEntity, err error) {
  if err = models.DB.First(&lang, id).Error; err != nil{

    return nil, err
  }

  return lang, nil
}

func UpdateLangById(id uint, edit dto.LangEdit) (lang *models.LangEntity, err error) {
  tx := models.DB.Begin()

  if err := tx.Error; err != nil {

      return nil, err
  }

  if err := models.DB.Table("langs").Where("id = ? And version = ?", id, edit.Version).First(&lang).Error; err != nil{
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
