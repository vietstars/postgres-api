package models

import (
  "gorm.io/gorm"
  "time"
)

type Quest struct {
  ID uint `json:"id" gorm:"primary_key"`  
  Title string `json:"title"`
  Description string `json:"description"`
  Reward int `json:"reward"`
  CategoryID uint `json:"category_id"`
  CreatedByID uint `gorm:"index;column:created_by" json:"created_by"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedByID uint `gorm:"index;column:updated_by" json:"updated_by"`
  UpdatedAt time.Time `json:"updated_at"`
  DeletedAt *time.Time `sql:"index" json:"-" swaggertype:"primitive,string"`
  Version uint `json:"version"`
}

func (q *Quest) BeforeSave(tx *gorm.DB) (err error) {
  tx.Statement.SetColumn("Version", q.Version+1)
  tx.Statement.SetColumn("CreatedAt", time.Now())

  return 
}

func (q *Quest) BeforeUpdate(tx *gorm.DB) (err error) {
  if tx.Statement.Changed() {
    tx.Statement.SetColumn("Version", q.Version+1)
    tx.Statement.SetColumn("UpdatedAt", time.Now())
  }

  return 
}

func NewQuest(title string, description string, reward int) (quest *Quest, err error){
  quest = &Quest{
    Title: title,
    Description: description,
    Reward: reward,
  }

  DB.Create(&quest)
   
  return quest, nil
}


