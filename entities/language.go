package entities

import (
  // "github.com/guregu/null"
  "gorm.io/gorm"
  "time"
)

type Lang struct {
  LangID uint `gorm:"primary_key;column:id" json:"lang_id"`
  Key string `gorm:"column:code" json:"key"`
  ValVN string `gorm:"column:pharse_vn" json:"val_vn"`
  ValEN string `gorm:"column:pharse_en" json:"val_en"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by_id"`
  CreatedAt time.Time `gorm:"column:created_at" json:"created_at_time"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by_id"`
  UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at_time"` //gorm:"default:current_timestamp"
  DeletedAt gorm.DeletedAt `json:"-"`
  Version uint `json:"version"`
}

type LangList []*Lang