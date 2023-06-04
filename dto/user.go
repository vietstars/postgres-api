package dto

import(
  "github.com/vietstars/postgres-api/models"
)

type UserNew struct {
  UserName string `json:"username" validate:"required"`
  Email string `json:"email" validate:"required"`
  Password string `json:"password" validate:"required"`
}

type UserSignIn struct {
  Email string `json:"email" validate:"required"`
  Password string `json:"password" validate:"required"`
}

type Auth struct {
  Info *models.UserEntity `json:"info"`
  Token string `json:"token"`
}
