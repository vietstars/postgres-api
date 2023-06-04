package seeders

import (
  "os"
  "gopkg.in/guregu/null.v4"

  "github.com/gocarina/gocsv"
  "github.com/vietstars/postgres-api/models"
)

type User struct {
  Name string `csv:"name"`
  Email string `csv:"email"`
  Password string `csv:"password"`
  Status string `csv:"status"`
  CreatedBy string `csv:"created_by"`
  CreatedAt string ` csv:"created_at"`
  UpdatedBy string `csv:"updated_by"`
  UpdatedAt string ` csv:"updated_at"`
  DeletedAt null.Time `csv:"deleted_at"`
  Version string `csv:"version"`
}

func SeedUser() {
  userFile, err := os.OpenFile("./files/users.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)

  if err != nil {
    panic(err)
  }
  defer userFile.Close()

  var users []*User

  err = gocsv.UnmarshalFile(userFile, &users)
  if err != nil {
    panic(err)
  }

  result := models.DB.Create(users)

  if result.Error != nil {
    panic(result.Error)
  }
}