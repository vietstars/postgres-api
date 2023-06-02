package repositories

import (
  "github.com/vietstars/postgres-api/models"
)

func NewLang(lg string, gr string, key string, val string) (lang *models.Lang, err error) {
  lang = &models.Lang{
    Locale: lg,
    Group: gr,
    Key: key,
    Val: val,
  }
  models.DB.Create(&lang)

  return lang, nil
}

func DelLangById(id uint, version uint, force bool) (result bool, err error) {
  var lang models.Lang

  tx := models.DB.Begin()

  if err := tx.Error; err != nil {

      return  false, err
  }

  if err := models.DB.Where("id = ? And version = ?", id, version).First(&lang).Error; err != nil{
    tx.Rollback()

    return false, err
  }

  if force {
    tx.Unscoped().Delete(&lang)
  } else {
    tx.Delete(&lang)
  }

  tx.Commit()

  return true, nil
}

func GetAllLang() (langs *models.LangList, err error) {
  if err = models.DB.Find(&langs).Error; err != nil{

    return nil, err
  }

  return langs, nil
}

func GetLangsByLocale(lg string) (langs *models.LangList, err error) {
  err = models.DB.Find(&langs, "locale = '_' OR locale = ?", lg).Error

  return langs, err
}

func GetLangById(id uint) (lang *models.Lang, err error) {
  if err = models.DB.First(&lang, id).Error; err != nil{

    return nil, err
  }

  return lang, nil
}

func UpdateLangById(id uint, version uint, lg string, gr string, key string, val string) (lang *models.Lang, err error) {

  tx := models.DB.Begin()

  if err := tx.Error; err != nil {

      return nil, err
  }

  if err := models.DB.Where("id = ? And version = ?", id, version).First(&lang).Error; err != nil{
    tx.Rollback()

    return nil, err
  }

  lang.Locale = lg
  lang.Group = gr 
  lang.Key = key 
  lang.Val = val

  tx.Save(&lang)
  tx.Commit()

  return lang, nil
}
