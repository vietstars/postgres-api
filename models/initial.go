package models

import (
  "fmt"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "os"
  "github.com/gocarina/gocsv"
  "github.com/vietstars/postgres-api/seeds"
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

  userImport, err := os.Open("users.csv")
  if err != nil {
    panic(err)
  }
  defer userImport.Close()

  var users []User

  err = gocsv.Unmarshal(userImport, &users)
  if err != nil {
    panic(err)
  }

  err = database.AutoMigrate(&User{}, &Store{}, &Lang{})

  if err != nil {
    panic(err)
  }

  result := DB.Create(users)
  if result.Error != nil {
    panic(result.Error)
  }

  DB = database
}
