package models

import (
  "gorm.io/gorm"
  "time"
)

type StoreEntity struct {
  ItemID uint `gorm:"primary_key;column:id" json:"item_id"`
  OwnerID uint `gorm:"index;column:owner" json:"owner_id"`
  ItemData string `gorm:"column:data" json:"item_data"`
  ItemDescription string `gorm:"column:description" json:"item_description"`
  ItemStatus uint8 `gorm:"column:status" json:"item_status"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by_id"`
  CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by_id"`
  UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
  DeletedAt gorm.DeletedAt `json:"-"`
  Version uint `gorm:"default:1" json:"version"`
  table string `gorm:"-"`
}

func (StoreEntity) TableName() string {
  return "stores"
}

func (s *StoreEntity) BeforeCreate(tx *gorm.DB) (err error) {
  tx.Statement.SetColumn("CreatedAt", time.Now())

  return 
}

func (s *StoreEntity) BeforeUpdate(tx *gorm.DB) (err error) {
    tx.Statement.SetColumn("Version", s.Version+1)
    tx.Statement.SetColumn("UpdatedAt", time.Now())

  return 
}
