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

func DelLangById(id uint, version uint, force bool) (err error, result bool) {
  var lang models.Lang

  tx := models.DB.Begin()

  if err := tx.Error; err != nil {

      return err, false
  }

  if err := models.DB.Where("id = ? And version = ?", id, version).First(&lang).Error; err != nil{
    tx.Rollback()

    return err, false
  }

  if force {
    tx.Unscoped().Delete(&lang)
  } else {
    tx.Delete(&lang)
  }

  tx.Commit()

  return nil, true
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
