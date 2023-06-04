package models

import (
  "fmt"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "os"
)

var DB *gorm.DB

func ConnectDatabase() {
  dsn := fmt.Sprintf(
    "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
  )
  database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
    panic("Failed to connect to database")
  }

  err = database.AutoMigrate(
    &UserEntity{table: "users"}, 
    &StoreEntity{table: "stores"}, 
    &LangEntity{table: "langs"},
  )

  if err != nil {
    panic(err)
  }

  DB = database
}
