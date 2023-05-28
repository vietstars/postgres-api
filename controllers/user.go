package controllers

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "net/http"

  "github.com/go-playground/validator/v10"
  "github.com/gorilla/mux"
  "github.com/vietstars/postgres-api/models"
  "github.com/vietstars/postgres-api/utils"
  "golang.org/x/crypto/bcrypt"
)

var userValidate *validator.Validate

type UserInput struct {
  Email string `json:"email" validate:"required"`
  Password string `json:"password" validate:"required"`
}

func GetAllUsers(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  var users []models.User
  models.DB.Find(&users)

  // for i,user := range users {
  //   if user.Version == 0 {
  //     users[i].Version = 1
  //   }
  // }

  json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")

  id := mux.Vars(r)["id"]
  var user models.User

  if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil{
    utils.RespondNotFound(w, 
      "User not found")
    return
  }

  json.NewEncoder(w).Encode(user)
}

func SignUp(w http.ResponseWriter, r *http.Request){
  var inputs UserInput 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &inputs)

  userValidate = validator.New()
  err := userValidate.Struct(inputs)

  if err != nil {
    utils.RespondBadRequest(w, 
      fmt.Sprintf("%+v\n", err))
    return 
  }

  user, err := models.NewUser(inputs.Email, inputs.Password)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(user) 
}

func SignIn(w http.ResponseWriter, r *http.Request){
  var inputs UserInput 

  body, _ := ioutil.ReadAll(r.Body)
  _ = json.Unmarshal(body, &inputs)

  userValidate = validator.New()
  err := userValidate.Struct(inputs)

  if err != nil {
    utils.RespondBadRequest(w, 
      fmt.Sprintf("%+v\n", err))
    return 
  }

  user, err := models.GetUserByEmail(inputs.Email)

  if err != nil {
    utils.RespondNotFound(w, 
      "User does not exist")
    return 
  }

  hashedPassword := user.Password
  err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputs.Password))

  if err != nil {
    utils.RespondBadRequest(w, 
      "Password does not match")
    return
  }

  token, err := utils.GenerateToken(user)
  if err != nil {
    utils.RespondServerError(w, 
      "Error while generating token")
    return
  }

  auth := models.Auth{user, token}

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(auth) 
}
