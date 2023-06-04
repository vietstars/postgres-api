package models

import (
  "gorm.io/gorm"
  "time"
)

type LangEntity struct {
  LangID uint `gorm:"primary_key;auto_increment;column:id" json:"lang_id"`
  Locale string `gorm:"index;column:locale" json:"lg"`
  Group string `gorm:"index;column:category" json:"group"`
  Key string `gorm:"column:code" json:"key"`
  Val string `gorm:"column:pharse" json:"val"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by_id"`
  CreatedAt time.Time `gorm:"column:created_at" json:"created_at_time"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by_id"`
  UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at_time"`
  DeletedAt gorm.DeletedAt `json:"-"`
  Version uint `gorm:"default:1" json:"version"`
  table string `gorm:"-"`
}

func (LangEntity) TableName() string {
  return "langs"
}

type LangListEntity []*LangEntity

func (l *LangEntity) BeforeCreate(tx *gorm.DB) (err error) {
  tx.Statement.SetColumn("CreatedAt", time.Now())

  return 
}

func (l *LangEntity) BeforeUpdate(tx *gorm.DB) (err error) {
  tx.Statement.SetColumn("Version", l.Version+1)
  tx.Statement.SetColumn("UpdatedAt", time.Now())

  return 
}